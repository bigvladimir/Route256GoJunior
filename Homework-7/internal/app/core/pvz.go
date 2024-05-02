package core

import (
	"context"

	"homework/internal/app/pvz/dto"
)

type pvzService interface {
	AddPvz(ctx context.Context, input dto.PvzInput) (int64, error)
	GetPvzByID(ctx context.Context, id int64) (dto.Pvz, error)
	UpdatePvz(ctx context.Context, input dto.Pvz) error
	DeletePvz(ctx context.Context, id int64) error
	ModifyPvz(Pvzctx context.Context, input dto.Pvz) (int64, error)
}

// AddPvz добавляет запись ПВЗ без указания id
func (s *Service) AddPvz(ctx context.Context, input dto.PvzInput) (int64, error) {
	return s.pvzService.AddPvz(ctx, input)
}

// GetPvzByID возвращает запись ПВЗ по id, если она существует
func (s *Service) GetPvzByID(ctx context.Context, id int64) (dto.Pvz, error) {
	return s.pvzService.GetPvzByID(ctx, id)
}

// UpdatePvz обновляет запись ПВЗ по id, если она существует
func (s *Service) UpdatePvz(ctx context.Context, input dto.Pvz) error {
	return s.pvzService.UpdatePvz(ctx, input)
}

// DeletePvz удаляет запись ПВЗ по id, если она существует
func (s *Service) DeletePvz(ctx context.Context, id int64) error {
	return s.pvzService.DeletePvz(ctx, id)
}

// ModifyPvz обновляет запись ПВЗ по id, если не находит, что обновить, то вставляет новую
func (s *Service) ModifyPvz(ctx context.Context, input dto.Pvz) (int64, error) {
	return s.pvzService.ModifyPvz(ctx, input)
}
