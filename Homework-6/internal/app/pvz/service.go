package pvz

import (
	"context"
	"errors"

	"homework/internal/app/pvz/dto"
)

type storage interface {
	Add(context.Context, dto.PvzInput) (int64, error)
	GetByID(context.Context, int64) (dto.Pvz, error)
	Update(context.Context, dto.Pvz) error
	Delete(context.Context, int64) error
	Modify(context.Context, dto.Pvz) (int64, error)
}

// Service provides functions for interacting with the pvz storage
type Service struct {
	stor storage
}

// NewService creates Service
func NewService(s storage) *Service {
	return &Service{stor: s}
}

// AddPvz добавляет запись без указания id
func (s *Service) AddPvz(ctx context.Context, input dto.PvzInput) (int64, error) {
	if input.Name == "" || input.Adress == "" || input.Contacts == "" {
		return -1, errors.New("Пустые поля")
	}
	return s.stor.Add(ctx, input)
}

// GetPvzByID возвращает строку по id, если она существует
func (s *Service) GetPvzByID(ctx context.Context, id int64) (dto.Pvz, error) {
	if id <= 0 {
		return dto.Pvz{}, errors.New("Некорректный id ПВЗ")
	}
	return s.stor.GetByID(ctx, id)
}

// UpdatePvz обновляет строку по id, если она существует
func (s *Service) UpdatePvz(ctx context.Context, input dto.Pvz) error {
	if input.ID <= 0 {
		return errors.New("Некорректный id ПВЗ")
	}
	if input.Name == "" || input.Adress == "" || input.Contacts == "" {
		return errors.New("Пустые поля")
	}
	return s.stor.Update(ctx, input)
}

// DeletePvz удаляет строку по id, если она существует
func (s *Service) DeletePvz(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("Некорректный id ПВЗ")
	}
	return s.stor.Delete(ctx, id)
}

// ModifyPvz обновляет данные по id, если не находит, что обновить, то вставляет новые
func (s *Service) ModifyPvz(ctx context.Context, input dto.Pvz) (int64, error) {
	if input.ID <= 0 {
		return -1, errors.New("Некорректный id ПВЗ")
	}
	if input.Name == "" || input.Adress == "" || input.Contacts == "" {
		return -1, errors.New("Пустые поля")
	}
	return s.stor.Modify(ctx, input)
}
