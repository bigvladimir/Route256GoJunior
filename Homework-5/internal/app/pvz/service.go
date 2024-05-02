package pvz

import (
	"context"
	"errors"
)

type Service struct {
	stor storage
}

type PvzService interface {
	AddPvz(ctx context.Context, input PvzInput) (int64, error)
	GetPvzByID(ctx context.Context, id int64) (Pvz, error)
	UpdatePvz(ctx context.Context, input Pvz) error
	DeletePvz(ctx context.Context, id int64) error
	ModifyPvz(Pvzctx context.Context, input Pvz) (int64, error)
}

// New инициализация Service
func NewService(s storage) *Service {
	return &Service{stor: s}
}

func (s *Service) AddPvz(ctx context.Context, input PvzInput) (int64, error) {
	if input.Name == "" || input.Adress == "" || input.Contacts == "" {
		return -1, errors.New("Пустые поля")
	}
	return s.stor.add(ctx, input)
}

func (s *Service) GetPvzByID(ctx context.Context, id int64) (Pvz, error) {
	if id <= 0 {
		return Pvz{}, errors.New("Некорректный id ПВЗ")
	}
	return s.stor.getByID(ctx, id)
}

func (s *Service) UpdatePvz(ctx context.Context, input Pvz) error {
	if input.ID <= 0 {
		return errors.New("Некорректный id ПВЗ")
	}
	if input.Name == "" || input.Adress == "" || input.Contacts == "" {
		return errors.New("Пустые поля")
	}
	return s.stor.update(ctx, input)
}

func (s *Service) DeletePvz(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("Некорректный id ПВЗ")
	}
	return s.stor.delete(ctx, id)
}

func (s *Service) ModifyPvz(ctx context.Context, input Pvz) (int64, error) {
	if input.ID <= 0 {
		return -1, errors.New("Некорректный id ПВЗ")
	}
	if input.Name == "" || input.Adress == "" || input.Contacts == "" {
		return -1, errors.New("Пустые поля")
	}
	return s.stor.modify(ctx, input)
}
