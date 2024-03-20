package service

import (
	"Homework-2/internal/model"
	"errors"
	"sort"
	"time"
)

type storage interface {
	Add(model.OrderInput) error
	Get(string, ...int) ([]model.Order, error)
	Update(model.Order) error
	Delete(int) error
}

type Service struct {
	s storage
}

// New инициализация Service
func New(s storage) Service {
	return Service{s: s}
}

func (s Service) add(input model.OrderInput) error {
	if input.OrderID <= 0 {
		return errors.New("Некорректный id заказа.")
	}
	if input.CustomerID <= 0 {
		return errors.New("Некорректный id клиента.")
	}
	if !input.StorageLastTime.After(time.Now()) {
		return errors.New("Некорректное время хранения заказа.")
	}

	return s.s.Add(input)
}

func (s Service) get(filter string, id ...int) ([]model.Order, error) {
	if len(filter) == 0 {
		err := errors.New("Пустой фильтр для Get().")
		return nil, err
	}
	if len(id) > 1 {
		err := errors.New("Слишком много аргументов для Get()")
		return nil, err
	}
	for _, j := range id {
		if j <= 0 {
			err := errors.New("Некорректный id для Get()")
			return nil, err
		}
	}

	return s.s.Get(filter, id...)
}

func (s Service) update(input model.Order) error {
	if input.OrderID <= 0 {
		return errors.New("Некорректный id заказа при обновлении.")
	}
	if input.CustomerID <= 0 {
		return errors.New("Некорректный id клиента при обновлении.")
	}
	if !input.IsCompleted && input.ArrivalTime.After(input.CompleteTime) {
		return errors.New("Некорректное время завершения заказа при обновлении.")
	}
	if input.ArrivalTime.After(input.StorageLastTime) {
		return errors.New("Некорректное время хранения заказа при обновлении.")
	}

	return s.s.Update(input)
}

func (s Service) delete(id int) error {
	if id <= 0 {
		return errors.New("Некорректный id заказа.")
	}

	return s.s.Delete(id)
}

func sortOrdersByArrival(orders []model.Order) {
	sort.Slice(orders, func(i, j int) bool {
		return orders[j].ArrivalTime.After(orders[i].ArrivalTime)
	})
}
