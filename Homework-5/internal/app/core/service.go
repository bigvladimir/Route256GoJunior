//go:generate mockgen -source ./service.go -destination=./mocks/service.go -package=mock_core
package core

import (
	"context"
	"homework/internal/app/orders"
	"homework/internal/app/pvz"
)

type Service struct {
	ordersService orders.OrdersService
	pvzService    pvz.PvzService
}

type CoreOps interface {
	TakeOrderFromCourier(ctx context.Context, order orders.OrderInput) error
	ReturnOrderToCourier(ctx context.Context, pvzID, orderID int) error
	GiveOrderToCustomer(ctx context.Context, pvzID int, orderID []int) error
	TakeRefundFromCustomer(ctx context.Context, pvzID, customerID, orderID int) error
	GetRefundList(ctx context.Context, pvzID, pageNum, pageSize int) ([]orders.Order, error)
	GetCustomerOrderList(ctx context.Context, pvzID, customerID, limit int, isInStock bool) ([]orders.Order, error)

	AddPvz(ctx context.Context, input pvz.PvzInput) (int64, error)
	GetPvzByID(ctx context.Context, id int64) (pvz.Pvz, error)
	UpdatePvz(ctx context.Context, input pvz.Pvz) error
	DeletePvz(ctx context.Context, id int64) error
	ModifyPvz(ctx context.Context, input pvz.Pvz) (int64, error)
}

func NewCoreService(ordersService orders.OrdersService, pvzService pvz.PvzService) *Service {
	return &Service{
		ordersService: ordersService,
		pvzService:    pvzService,
	}
}
