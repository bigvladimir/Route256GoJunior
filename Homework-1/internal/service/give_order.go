package service

import (
	"Homework-1/internal/model"
	"errors"
	"fmt"
	"log"
	"time"
)

// GiveOrderToCustomer обрабатывает выдачу заказов одному клиенту
func (s Service) GiveOrderToCustomer(orderID []int) error {
	var orders []model.Order
	for _, id := range orderID {
		order, err := s.get("orderID", id)
		if err != nil {
			return fmt.Errorf("Не удалось получить информацию о заказах: %w", err)
		}
		if len(order) == 0 {
			return errors.New("Не все заказы найдены")
		}
		if order[0].IsCompleted && !order[0].IsRefunded {
			return errors.New("Попытка забрать уже отданный заказ.")
		}
		if order[0].IsRefunded {
			return errors.New("Попытка забрать возвращенный заказ.")
		}
		if time.Now().After(order[0].StorageLastTime) {
			return errors.New("Не у всех заказов действительный срок хранения.")
		}

		orders = append(orders, order[0])
	}

	checkSameCustomer := orders[0].CustomerID
	for i := 1; i < len(orders); i++ {
		if orders[i].CustomerID != checkSameCustomer {
			return errors.New("Не все заказы принаддежат одному клиенту")
		}
	}

	for _, order := range orders {
		order.IsCompleted = true
		err := s.update(order)
		if err != nil {
			return fmt.Errorf("Не удалось обновить информацию о заказе в базе данных: %w", err)
		}
		log.Println("Выдан заказ ID", order.OrderID)
	}

	return nil
}
