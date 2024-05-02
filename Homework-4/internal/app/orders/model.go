package orders

import "time"

type orderRow struct {
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

func (t *orderRow) mapToDTO() Order {
	return Order{
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
