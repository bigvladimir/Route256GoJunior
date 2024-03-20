package service

import (
	"errors"
	"fmt"
	"time"
)

// TakeRefundFromCustomer обрабатывает возврат заказа клиентом
func (s Service) TakeRefundFromCustomer(customerID, orderID int) error {
	orders, err := s.get("orderID", orderID)
	if err != nil {
		return fmt.Errorf("Не удалось получить информацию о заказе: %w", err)
	}
	if len(orders) == 0 {
		return errors.New("Заказ не выдавался в этом ПВЗ.")
	}
	order := orders[0]
	if customerID != order.CustomerID {
		return errors.New("У заказа другой владелец.")
	}
	if !order.IsCompleted {
		return errors.New("Заказ найден, но он не выдавался никому.")
	}
	if order.IsRefunded {
		return errors.New("Заказ найден, но уже вернули.")
	}

	lastTimeForRefund := time.Date(
		order.StorageLastTime.Year(),
		order.StorageLastTime.Month(),
		order.StorageLastTime.Day()+2,
		23,
		59,
		0,
		0,
		order.StorageLastTime.Location(),
	)
	if time.Now().After(lastTimeForRefund) {
		return errors.New("Срок возврата истёк.")
	}

	order.IsRefunded = true
	err = s.update(order)
	if err != nil {
		return fmt.Errorf("Не удалось обновить информацию о заказе в базе данных: %w", err)
	}

	return nil
}
