package service

import (
	"fmt"
	"homework/internal/interactive/model"
	"time"
)

// TakeOrderFromCourier обратывает принятие заказа от курьера
func (s Service) TakeOrderFromCourier(order model.OrderInput) error {
	order.StorageLastTime = time.Date(
		order.StorageLastTime.Year(),
		order.StorageLastTime.Month(),
		order.StorageLastTime.Day(),
		23,
		59,
		0,
		0,
		time.Now().Location(),
	)

	err := s.add(order)
	if err != nil {
		return fmt.Errorf("Не удалось добавить заказ в базу данных: %w", err)
	}

	return nil
}
