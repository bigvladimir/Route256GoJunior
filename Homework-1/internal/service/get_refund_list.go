package service

import (
	"Homework-1/internal/model"
	"errors"
	"fmt"
)

// GetRefundList возвращает страницу возвращенных заказов в виде слайса
// pageNum int - номер страницы, pageSize int - размер страницы
func (s Service) GetRefundList(pageNum, pageSize int) ([]model.Order, error) {
	orders, err := s.get("isRefunded")
	if err != nil {
		err = fmt.Errorf("Не удалось получить информацию о заказах: %w", err)
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
