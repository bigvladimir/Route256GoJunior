package core

import (
	"context"
	"homework/internal/app/orders"
	"homework/internal/app/pvz"
)

type Service struct {
	ordersService OrdersService
	pvzService    PvzService
}

type OrdersService interface {
	TakeOrderFromCourier(ctx context.Context, order orders.OrderInput) error
	ReturnOrderToCourier(ctx context.Context, pvzID, orderID int) error
	GiveOrderToCustomer(ctx context.Context, pvzID int, orderID []int) error
	TakeRefundFromCustomer(ctx context.Context, pvzID, customerID, orderID int) error
	GetRefundList(ctx context.Context, pvzID, pageNum, pageSize int) ([]orders.Order, error)
	GetCustomerOrderList(ctx context.Context, pvzID, customerID, limit int, isInStock bool) ([]orders.Order, error)
}

type PvzService interface {
	AddPvz(ctx context.Context, input pvz.PvzInput) (int64, error)
	GetPvzByID(ctx context.Context, id int64) (pvz.Pvz, error)
	UpdatePvz(ctx context.Context, input pvz.Pvz) error
	DeletePvz(ctx context.Context, id int64) error
	ModifyPvz(Pvzctx context.Context, input pvz.Pvz) (int64, error)
}

func NewCoreService(ordersService OrdersService, pvzService PvzService) *Service {
	return &Service{
		ordersService: ordersService,
		pvzService:    pvzService,
	}
}
