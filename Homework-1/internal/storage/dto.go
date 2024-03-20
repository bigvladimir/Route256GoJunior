package storage

import "time"

type orderDTO struct {
	OrderID         int
	CustomerID      int
	StorageLastTime time.Time

	IsCompleted  bool
	CompleteTime time.Time

	IsRefunded bool

	ArrivalTime time.Time
}
