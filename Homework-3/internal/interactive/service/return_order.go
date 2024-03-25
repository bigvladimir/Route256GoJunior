package service

import (
	"errors"
	"fmt"
	"time"
)

// ReturnOrderToCourier обратывает возврат заказа курьеру,
// если у заказа закончился срок хранения,
// позволяет отдать курьеру заказ возвращенный клиентом независимо от срока хранения
func (s Service) ReturnOrderToCourier(orderID int) error {
	orders, err := s.get("orderID", orderID)
	if err != nil {
		return fmt.Errorf("Не удалось получить информацию о заказе: %w", err)
	}
	if len(orders) == 0 {
		return errors.New("Заказ не найден.")
	}

	order := orders[0]
	if order.IsCompleted && !order.IsRefunded {
		return errors.New("Заказ уже был выдан клиенту.")
	}
	if order.StorageLastTime.After(time.Now()) && !order.IsRefunded {
		return errors.New("У заказа ещё не вышел срок хранения.")
	}

	err = s.delete(order.OrderID)
	if err != nil {
		return fmt.Errorf("Не удалось удалить заказ из базы данных: %w", err)
	}

	return nil
}
