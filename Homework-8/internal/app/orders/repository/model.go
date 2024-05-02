package repository

import (
	"time"

	"homework/internal/app/orders/dto"
)

// OrderRow is the package internal representation of the order
type OrderRow struct {
	OrderID         int       `db:"order_id"`
	PvzID           int       `db:"pvz_id"`
	CustomerID      int       `db:"customer_id"`
	StorageLastTime time.Time `db:"storage_last_time"`
	IsCompleted     bool      `db:"is_completed"`
	CompleteTime    time.Time `db:"complete_time"`
	IsRefunded      bool      `db:"is_refunded"`
	ArrivalTime     time.Time `db:"arrival_time"`
	PackageType     string    `db:"package_type"`
	Weight          float64   `db:"weight"`
	Price           int       `db:"price"`
}

// MapFromDTO converts external representation to OrderRow
func (t *OrderRow) MapFromDTO(d dto.OrderInput) {
	t.OrderID = d.OrderID
	t.PvzID = d.PvzID
	t.CustomerID = d.CustomerID
	t.StorageLastTime = d.StorageLastTime
	t.IsCompleted = false
	t.CompleteTime = dto.ZeroTime
	t.IsRefunded = false
	t.ArrivalTime = time.Now()
	t.PackageType = d.PackageType
	t.Weight = d.Weight
	t.Price = d.Price
}

// MapToDTO converts OrderRow to external representation
func (t *OrderRow) MapToDTO() dto.Order {
	return dto.Order{
		OrderID:         t.OrderID,
		PvzID:           t.PvzID,
		CustomerID:      t.CustomerID,
		StorageLastTime: t.StorageLastTime,
		IsCompleted:     t.IsCompleted,
		CompleteTime:    t.CompleteTime,
		IsRefunded:      t.IsRefunded,
		ArrivalTime:     t.ArrivalTime,
		PackageType:     t.PackageType,
		Weight:          t.Weight,
		Price:           t.Price,
	}
}
