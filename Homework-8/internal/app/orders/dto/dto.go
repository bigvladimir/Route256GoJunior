package dto

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
