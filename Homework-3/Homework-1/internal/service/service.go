package service

import (
	"Homework-1/internal/model"
	"errors"
	"log"
	"sort"
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
		return err
	}

	return nil
}

// ReturnOrderToCourier обратывает возврат заказа курьеру,
// если у заказа закончился срок хранения,
// позволяет отдать курьеру заказ возвращенный клиентом независимо от срока хранения
func (s Service) ReturnOrderToCourier(orderID int) error {
	orders, err := s.get("orderID", orderID)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// GiveOrderToCustomer обрабатывает выдачу заказов одному клиенту
func (s Service) GiveOrderToCustomer(orderID []int) error {
	var orders []model.Order
	for _, id := range orderID {
		order, err := s.get("orderID", id)
		if err != nil {
			return err
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
			return err
		}
		log.Println("Выдан заказ ID", order.OrderID)
	}

	return nil
}

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

// TakeRefundFromCustomer обрабатывает возврат заказа клиентом
func (s Service) TakeRefundFromCustomer(customerID, orderID int) error {
	orders, err := s.get("orderID", orderID)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// GetRefundList возвращает страницу возвращенных заказов в виде слайса
// pageNum int - номер страницы, pageSize int - размер страницы
func (s Service) GetRefundList(pageNum, pageSize int) ([]model.Order, error) {
	orders, err := s.get("isRefunded")
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		err := errors.New("Возвраты не найдены.")
		return nil, err
	}
	if pageSize*(pageNum-1) >= len(orders) {
		err := errors.New("Возвратов меньше чем запрашиваемая страница.")
		return nil, err
	}

	limit := pageSize * (pageNum)
	if cap(orders) < limit {
		limit = cap(orders)
	}
	return orders[pageSize*(pageNum-1) : limit], nil
}

func sortOrdersByArrival(orders []model.Order) {
	sort.Slice(orders, func(i, j int) bool {
		return orders[j].ArrivalTime.After(orders[i].ArrivalTime)
	})
}
