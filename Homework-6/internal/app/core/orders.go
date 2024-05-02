package core

import (
	"context"

	"homework/internal/app/orders"
)

type ordersService interface {
	TakeOrderFromCourier(ctx context.Context, order orders.OrderInput) error
	ReturnOrderToCourier(ctx context.Context, pvzID, orderID int) error
	GiveOrderToCustomer(ctx context.Context, pvzID int, orderID []int) error
	TakeRefundFromCustomer(ctx context.Context, pvzID, customerID, orderID int) error
	GetRefundList(ctx context.Context, pvzID, pageNum, pageSize int) ([]orders.Order, error)
	GetCustomerOrderList(ctx context.Context, pvzID, customerID, limit int, isInStock bool) ([]orders.Order, error)
}

func (s *Service) TakeOrderFromCourier(ctx context.Context, order orders.OrderInput) error {
	return s.ordersService.TakeOrderFromCourier(ctx, order)
}

func (s *Service) ReturnOrderToCourier(ctx context.Context, pvzID, orderID int) error {
	return s.ordersService.ReturnOrderToCourier(ctx, pvzID, orderID)
}

func (s *Service) GiveOrderToCustomer(ctx context.Context, pvzID int, orderID []int) error {
	return s.ordersService.GiveOrderToCustomer(ctx, pvzID, orderID)
}

func (s *Service) TakeRefundFromCustomer(ctx context.Context, pvzID, customerID, orderID int) error {
	return s.ordersService.TakeRefundFromCustomer(ctx, pvzID, customerID, orderID)
}

func (s *Service) GetRefundList(ctx context.Context, pvzID, pageNum, pageSize int) ([]orders.Order, error) {
	return s.ordersService.GetRefundList(ctx, pvzID, pageNum, pageSize)
}

func (s *Service) GetCustomerOrderList(ctx context.Context, pvzID, customerID, limit int, isInStock bool) ([]orders.Order, error) {
	return s.ordersService.GetCustomerOrderList(ctx, pvzID, customerID, limit, isInStock)
}
