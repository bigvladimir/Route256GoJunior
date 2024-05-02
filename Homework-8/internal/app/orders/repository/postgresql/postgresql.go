package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"

	"homework/internal/app/orders/dto"
	orders_errors "homework/internal/app/orders/errors"
	"homework/internal/app/orders/repository"
	"homework/internal/pkg/db/transaction_manager"
)

// не совсем понимаю как можно убрать зависимость от этого интерфейса
type queryEngineManager interface {
	GetQueryEngine(ctx context.Context) transaction_manager.DbOps
}

// OrderStorage allows you to use order database operations
type OrderStorage struct {
	qe queryEngineManager
}

// NewOrderStorage creates OrderStorage
func NewOrderStorage(queryEngine queryEngineManager) *OrderStorage {
	return &OrderStorage{qe: queryEngine}
}

// Add добавляет запись
func (s *OrderStorage) Add(ctx context.Context, input dto.OrderInput) error {
	var order repository.OrderRow
	order.MapFromDTO(input)
	_, err := s.qe.GetQueryEngine(ctx).Exec(
		ctx, `INSERT INTO
		      orders(order_id, pvz_id, customer_id, storage_last_time, is_completed, complete_time, is_refunded, arrival_time, package_type, weight, price)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
		order.OrderID, order.PvzID, order.CustomerID, order.StorageLastTime,
		order.IsCompleted, order.CompleteTime, order.IsRefunded, order.ArrivalTime,
		order.PackageType, order.Weight, order.Price,
	)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		// код ошибки нарушения внешнего ключа
		if ok && pgErr.Code == "23503" {
			// TODO пока в таблице один внешний ключ, то можно сделать так, но нужно более точно определять ошибку
			return errors.New("Такого ПВЗ не существует")
		}
		return fmt.Errorf("Ошибка при записи данных: %w", err)
	}
	return nil
}

// Get возвращает запись по фильтру
func (s *OrderStorage) Get(ctx context.Context, filter string, id ...int) ([]dto.Order, error) {
	var row []repository.OrderRow
	var err error

	switch filter {
	case "orderID":
		if len(id) != 1 {
			return nil, errors.New("Неправильное количество аргументов для фильтра orderID в Get()")
		}
		err = s.qe.GetQueryEngine(ctx).Select(ctx, &row, "SELECT * FROM orders WHERE order_id=$1;", id[0])
	case "customerID":
		if len(id) != 2 {
			return nil, errors.New("Неправильное количество аргументов для фильтра customerID в Get()")
		}
		err = s.qe.GetQueryEngine(ctx).Select(ctx, &row, "SELECT * FROM orders WHERE pvz_id = $1 AND customer_id=$2;", id[0], id[1])
	case "isRefunded":
		if len(id) != 2 {
			return nil, errors.New("Неправильное количество аргументов для фильтра isRefunded в Get()")
		}
		err = s.qe.GetQueryEngine(ctx).Select(ctx, &row, "SELECT * FROM orders WHERE is_refunded IS TRUE AND pvz_id = $1;", id[0])
	default:
		return nil, errors.New("Непредусмотренный фильтр в Get()")
	}

	if err != nil {
		return nil, fmt.Errorf("Не удалось получить данные по фильтру: %w", err)
	}

	var orders []dto.Order

	for _, order := range row {
		orders = append(orders, order.MapToDTO())
	}

	return orders, nil
}

// Update обновляет строку по id, если она существует
func (s *OrderStorage) Update(ctx context.Context, id int, column, value string) error {
	commandTag, err := s.qe.GetQueryEngine(ctx).Exec(
		ctx, "UPDATE orders SET $1 = $2 WHERE id = $3;", column, value, id,
	)
	if err != nil {
		err = fmt.Errorf("Ошибка при обновлении данных: %w", err)
		return err
	}
	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return orders_errors.ErrNotFound
	}
	return nil
}

// Delete удаляет строку по id, если она существует
func (s *OrderStorage) Delete(ctx context.Context, id int) error {
	commandTag, err := s.qe.GetQueryEngine(ctx).Exec(ctx, "DELETE FROM orders WHERE id = $1;", id)
	if err != nil {
		err = fmt.Errorf("Ошибка при удалении данных: %w", err)
		return err
	}
	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return orders_errors.ErrNotFound
	}
	return nil
}
