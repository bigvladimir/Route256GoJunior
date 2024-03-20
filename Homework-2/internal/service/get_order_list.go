package service

import (
	"Homework-2/internal/model"
	"errors"
	"fmt"
)

// GetCustomerOrderList возвращает слайс заказов по ID клиента,
// limit int устанавливает максимальное количество возвращаемых заказов,
// если limit = 0, то ограничения нет
// isInStock bool устанавливает необходимость проверки наличия заказ в пункте,
// в том числе возвращенные
func (s Service) GetCustomerOrderList(customerID, limit int, isInStock bool) ([]model.Order, error) {
	if limit < 0 {
		err := errors.New("Отрицательное значение ограничения количества заказов.")
		return nil, err
	}
	orders, err := s.get("customerID", customerID)
	if err != nil {
		err = fmt.Errorf("Не удалось получить информацию о заказах: %w", err)
		return nil, err
	}
	if len(orders) == 0 {
		err := errors.New("Заказы не найдены.")
		return nil, err
	}

	var newOrders []model.Order
	if isInStock {
		for _, order := range orders {
			if !order.IsCompleted || order.IsRefunded {
				newOrders = append(newOrders, order)
			}
		}
		orders = newOrders
	}

	sortOrdersByArrival(orders)
	if limit == 0 || limit > len(orders) {
		limit = len(orders)
	}

	return orders[:limit], nil
}
