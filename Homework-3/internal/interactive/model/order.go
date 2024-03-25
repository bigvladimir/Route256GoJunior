package model

import "time"

type Order struct {
	OrderID         int
	CustomerID      int
	StorageLastTime time.Time

	IsCompleted     bool
	CompleteTime    time.Time

	IsRefunded      bool

	ArrivalTime     time.Time
}

type OrderInput struct {
	OrderID         int
	CustomerID      int
	StorageLastTime time.Time
}
