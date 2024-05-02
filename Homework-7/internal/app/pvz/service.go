package pvz

import (
	"context"
	"fmt"
	"log"

	"homework/internal/app/pvz/dto"
	"homework/internal/app/pvz/validation"
	"homework/internal/pkg/cacheupdater"
)

type storage interface {
	Add(context.Context, dto.PvzInput) (int64, error)
	GetByID(context.Context, int64) (dto.Pvz, error)
	Update(context.Context, dto.Pvz) error
	Delete(context.Context, int64) error
	Modify(context.Context, dto.Pvz) (int64, error)
}

type transactionManager interface {
	RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error
}

type inMemoryCacheOps interface {
	GetPvzType(id int64) (dto.Pvz, bool)
	SetPvzType(dto.Pvz)
}

type cacheUpdaterOps interface {
	SendMessage(message cacheupdater.Message) error
}

// Service provides functions for validating and interacting with the pvz storage
type Service struct {
	stor          storage
	txManager     transactionManager
	inMemoryCache inMemoryCacheOps
	cacheUpdater  cacheUpdaterOps
}

// NewService creates Service
func NewService(s storage, tx transactionManager, inMemoryCache inMemoryCacheOps, cacheUpdater cacheUpdaterOps) *Service {
	return &Service{
		stor:          s,
		txManager:     tx,
		inMemoryCache: inMemoryCache,
		cacheUpdater:  cacheUpdater,
	}
}

// AddPvz добавляет запись без указания id, также валидирует значения
func (s *Service) AddPvz(ctx context.Context, input dto.PvzInput) (int64, error) {
	if err := validation.ValidatePvzInput(input); err != nil {
		return -1, fmt.Errorf("validation.ValidatePvzInput: %w", err)
	}

	return s.stor.Add(ctx, input)
}

// GetPvzByID возвращает строку по id, если она существует, также валидирует значения
func (s *Service) GetPvzByID(ctx context.Context, id int64) (dto.Pvz, error) {
	if err := validation.ValidatePvzID(id); err != nil {
		return dto.Pvz{}, fmt.Errorf("validation.ValidatePvzID: %w", err)
	}

	if pvzObj, ok := s.inMemoryCache.GetPvzType(id); ok {
		return pvzObj, nil
	}

	pvzObj, err := s.stor.GetByID(ctx, id)
	if err != nil {
		return dto.Pvz{}, err
	}

	go s.inMemoryCache.SetPvzType(pvzObj)

	return pvzObj, nil
}

// UpdatePvz обновляет строку по id, если она существует, также валидирует значения
func (s *Service) UpdatePvz(ctx context.Context, input dto.Pvz) error {
	if err := validation.ValidatePvz(input); err != nil {
		return fmt.Errorf("validation.ValidatePvz: %w", err)
	}

	if pvzObj, ok := s.inMemoryCache.GetPvzType(input.ID); ok && input == pvzObj {
		return nil
	}

	if err := s.txManager.RunSerializable(ctx,
		func(ctxTX context.Context) error {
			return s.stor.Update(ctxTX, input)
		},
	); err != nil {
		return err
	}

	go func() {
		if err := s.cacheUpdater.SendMessage(cacheupdater.Message{ID: input.ID}); err != nil {
			log.Println("cacheUpdater.SendMessage:", err)
		}
	}()

	return nil
}

// DeletePvz удаляет строку по id, если она существует, также валидирует значения
func (s *Service) DeletePvz(ctx context.Context, id int64) error {
	if err := validation.ValidatePvzID(id); err != nil {
		return fmt.Errorf("validation.ValidatePvzID: %w", err)
	}

	if err := s.txManager.RunSerializable(ctx,
		func(ctxTX context.Context) error {
			return s.stor.Delete(ctxTX, id)
		},
	); err != nil {
		return err
	}

	go func() {
		if err := s.cacheUpdater.SendMessage(cacheupdater.Message{ID: id}); err != nil {
			log.Println("cacheUpdater.SendMessage:", err)
		}
	}()

	return nil
}

// ModifyPvz обновляет данные по id, если не находит, что обновить, то вставляет новые, также валидирует значения
func (s *Service) ModifyPvz(ctx context.Context, input dto.Pvz) (int64, error) {
	if err := validation.ValidatePvz(input); err != nil {
		return 0, fmt.Errorf("validation.ValidatePvz: %w", err)
	}

	if pvzObj, ok := s.inMemoryCache.GetPvzType(input.ID); ok && input == pvzObj {
		return 0, nil
	}

	var returningID *int64
	if err := s.txManager.RunSerializable(ctx,
		func(ctxTX context.Context) error {
			value, err := s.stor.Modify(ctxTX, input)
			returningID = &value
			return err
		},
	); err != nil {
		return 0, err
	}

	go func() {
		if err := s.cacheUpdater.SendMessage(cacheupdater.Message{ID: input.ID}); err != nil {
			log.Println("cacheUpdater.SendMessage:", err)
		}
	}()

	return *returningID, nil
}
