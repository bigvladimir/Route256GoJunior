package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type dbOps interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetPool(_ context.Context) *pgxpool.Pool
}

type OrderStorage struct {
	db dbOps
}

// New инициализирует Storage
func NewOrderStorage(database dbOps) *OrderStorage {
	return &OrderStorage{db: database}
}

// Add добавляет новую запись, возвращает ошибку есои запись существует
func (s *OrderStorage) add(ctx context.Context, input OrderInput) error {
	order := input.mapToModel()
	// выглядит ужасно, но не знаю как сделать лучше
	_, err := s.db.Exec(
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

// Get возвращает слайс заказов по фильтру
func (s *OrderStorage) get(ctx context.Context, filter string, id ...int) ([]Order, error) {
	var row []orderRow
	var err error
	// хотелось бы избежать кучи проверок типов, поэтому я сделаю запросы прямо внутри свича.
	// а не буду его составлять внутри
	switch filter {
	case "orderID":
		if len(id) != 1 {
			return nil, errors.New("Неправильное количество аргументов для фильтра orderID в Get()")
		}
		err = s.db.Select(ctx, &row, "SELECT * FROM orders WHERE order_id=$1;", id[0])
	case "customerID":
		if len(id) != 1 {
			return nil, errors.New("Неправильное количество аргументов для фильтра customerID в Get()")
		}
		err = s.db.Select(ctx, &row, "SELECT * FROM orders WHERE pvz_id = $1 AND customer_id=$2;", id[0], id[1])
	case "isRefunded":
		if len(id) != 2 {
			return nil, errors.New("Неправильное количество аргументов для фильтра isRefunded в Get()")
		}
		err = s.db.Select(ctx, &row, "SELECT * FROM orders WHERE is_refunded IS TRUE AND pvz_id = $1;", id[0])
	default:
		return nil, errors.New("Непредусмотренный фильтр в Get()")
	}

	if err != nil {
		return nil, fmt.Errorf("Не удалось получить данные по фильтру: %w", err)
	}

	var orders []Order

	for _, order := range row {
		orders = append(orders, order.mapToDTO())
	}

	return orders, nil
}

// Update обновляет одну колонку в записи, если она присутствовала, иначе возвращает ошибку
// TODO в коде не требуется фнукционал полного обновления информации о заказе, а только
// по одной колонке, но наверное надо переделать в обновлении по нескольким столбцам или всей записи
// поскольку хочу переделать пока не буду вешать кучу проверок на фильтр перед запросом
func (s *OrderStorage) update(ctx context.Context, id int, column, value string) error {
	commandTag, err := s.db.Exec(
		ctx, "UPDATE orders SET $1 = $2 WHERE id = $3;", column, value, id,
	)
	if err != nil {
		err = fmt.Errorf("Ошибка при обновлении данных: %w", err)
		return err
	}
	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete удаляет запись из базы, если она присутствовала, иначе возвращает ошибку
func (s *OrderStorage) delete(ctx context.Context, id int) error {
	commandTag, err := s.db.Exec(ctx, "DELETE FROM orders WHERE id = $1;", id)
	if err != nil {
		err = fmt.Errorf("Ошибка при удалении данных: %w", err)
		return err
	}
	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
