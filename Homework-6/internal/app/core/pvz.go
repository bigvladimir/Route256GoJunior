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

func (s *Service) AddPvz(ctx context.Context, input dto.PvzInput) (int64, error) {
	return s.pvzService.AddPvz(ctx, input)
}

func (s *Service) GetPvzByID(ctx context.Context, id int64) (dto.Pvz, error) {
	return s.pvzService.GetPvzByID(ctx, id)
}

func (s *Service) UpdatePvz(ctx context.Context, input dto.Pvz) error {
	return s.pvzService.UpdatePvz(ctx, input)
}

func (s *Service) DeletePvz(ctx context.Context, id int64) error {
	return s.pvzService.DeletePvz(ctx, id)
}

func (s *Service) ModifyPvz(ctx context.Context, input dto.Pvz) (int64, error) {
	return s.pvzService.ModifyPvz(ctx, input)
}
