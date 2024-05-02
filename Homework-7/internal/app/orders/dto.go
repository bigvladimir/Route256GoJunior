package orders

import (
	"time"
)

// ZeroTime is the default value for time.Time
var ZeroTime = time.Time{}

// Order is a external representation of the order
type Order struct {
	OrderID         int
	PvzID           int
	CustomerID      int
	StorageLastTime time.Time
	IsCompleted     bool
	CompleteTime    time.Time
	IsRefunded      bool
	ArrivalTime     time.Time
	PackageType     string
	Weight          float64
	Price           int
}

// OrderInput is a short external representation of the order
type OrderInput struct {
	OrderID         int
	PvzID           int
	CustomerID      int
	StorageLastTime time.Time
	PackageType     string
	Weight          float64
	Price           int
}

func (t *OrderInput) mapToModel() orderRow {
	return orderRow{
		OrderID:         t.OrderID,
		PvzID:           t.PvzID,
		CustomerID:      t.CustomerID,
		StorageLastTime: t.StorageLastTime,
		IsCompleted:     false,
		CompleteTime:    ZeroTime,
		IsRefunded:      false,
		ArrivalTime:     time.Now(),
		PackageType:     t.PackageType,
		Weight:          t.Weight,
		Price:           t.Price,
	}
}
