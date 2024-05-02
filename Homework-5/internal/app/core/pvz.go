package core

import (
	"context"
	"homework/internal/app/pvz"
)

func (s *Service) AddPvz(ctx context.Context, input pvz.PvzInput) (int64, error) {
	return s.pvzService.AddPvz(ctx, input)
}

func (s *Service) GetPvzByID(ctx context.Context, id int64) (pvz.Pvz, error) {
	return s.pvzService.GetPvzByID(ctx, id)
}

func (s *Service) UpdatePvz(ctx context.Context, input pvz.Pvz) error {
	return s.pvzService.UpdatePvz(ctx, input)
}

func (s *Service) DeletePvz(ctx context.Context, id int64) error {
	return s.pvzService.DeletePvz(ctx, id)
}

func (s *Service) ModifyPvz(ctx context.Context, input pvz.Pvz) (int64, error) {
	return s.pvzService.ModifyPvz(ctx, input)
}
