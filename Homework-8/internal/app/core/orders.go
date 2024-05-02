package core

import (
	"context"

	"homework/internal/app/orders/dto"
)

type ordersService interface {
	TakeOrderFromCourier(ctx context.Context, order dto.OrderInput) error
	ReturnOrderToCourier(ctx context.Context, pvzID, orderID int) error
	GiveOrderToCustomer(ctx context.Context, pvzID, customerID, orderID int) error
	TakeRefundFromCustomer(ctx context.Context, pvzID, customerID, orderID int) error
	GetRefundList(ctx context.Context, pvzID, pageNum, pageSize int) ([]dto.Order, error)
	GetCustomerOrderList(ctx context.Context, pvzID, customerID, limit int, isInStock bool) ([]dto.Order, error)
}

// TakeOrderFromCourier обратывает принятие заказа от курьера
func (s *Service) TakeOrderFromCourier(ctx context.Context, order dto.OrderInput) error {
	return s.ordersService.TakeOrderFromCourier(ctx, order)
}

// ReturnOrderToCourier обратывает возврат заказа курьеру
func (s *Service) ReturnOrderToCourier(ctx context.Context, pvzID, orderID int) error {
	return s.ordersService.ReturnOrderToCourier(ctx, pvzID, orderID)
}

// GiveOrderToCustomer обрабатывает выдачу заказов одному клиенту
func (s *Service) GiveOrderToCustomer(ctx context.Context, pvzID, customerID, orderID int) error {
	return s.ordersService.GiveOrderToCustomer(ctx, pvzID, customerID, orderID)
}

// TakeRefundFromCustomer обрабатывает возврат заказа клиентом
func (s *Service) TakeRefundFromCustomer(ctx context.Context, pvzID, customerID, orderID int) error {
	return s.ordersService.TakeRefundFromCustomer(ctx, pvzID, customerID, orderID)
}

// GetRefundList возвращает страницу возвращенных заказов в этом пвз в виде слайса
// pageNum int - номер страницы, pageSize int - размер страницы
func (s *Service) GetRefundList(ctx context.Context, pvzID, pageNum, pageSize int) ([]dto.Order, error) {
	return s.ordersService.GetRefundList(ctx, pvzID, pageNum, pageSize)
}

// GetCustomerOrderList возвращает слайс заказов по ID клиента в этом ПВЗ,
// limit int устанавливает максимальное количество возвращаемых заказов,
// если limit = 0, то ограничения нет
// isInStock bool устанавливает необходимость проверки наличия заказ в пункте,
// в том числе возвращенные
func (s *Service) GetCustomerOrderList(ctx context.Context, pvzID, customerID, limit int, isInStock bool) ([]dto.Order, error) {
	return s.ordersService.GetCustomerOrderList(ctx, pvzID, customerID, limit, isInStock)
}
