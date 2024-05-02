package grpcserver

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	orders_errors "homework/internal/app/orders/errors"
	"homework/internal/pkg/pb"
)

// TakeOrderFromCourier обратывает принятие заказа от курьера
func (s *Server) TakeOrderFromCourier(ctx context.Context, req *pb.OrderInfo) (*emptypb.Empty, error) {
	spanCtx, span := s.tracer.Start(ctx, "TakeOrderFromCourier")
	defer span.End()

	var orderInfo orderInfoModel
	orderInfo.mapFromProto(req)

	if err := orderInfo.validate(); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "orderInfoModel.validate: %v", err)
	}
	if err := s.service.TakeOrderFromCourier(spanCtx, orderInfo.mapToDTO()); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "TakeOrderFromCourier: %v", err)
	}

	s.metrics.IncOrdersBalance()

	return &emptypb.Empty{}, nil
}

// ReturnOrderToCourier обратывает возврат заказа курьеру
func (s *Server) ReturnOrderToCourier(ctx context.Context, req *pb.OrderIdentifier) (*emptypb.Empty, error) {
	spanCtx, span := s.tracer.Start(ctx, "ReturnOrderToCourier")
	defer span.End()

	var orderIdentifier orderIdentifierModel
	orderIdentifier.mapFromProto(req)

	if err := orderIdentifier.validate(); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "orderIdentifierModel.validate: %v", err)
	}
	if err := s.service.ReturnOrderToCourier(spanCtx, orderIdentifier.pvzID, orderIdentifier.orderID); err != nil {
		if errors.Is(err, orders_errors.ErrNotFound) {
			return &emptypb.Empty{}, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "ReturnOrderToCourier: %v", err)
	}

	s.metrics.IncReturnedOrders()
	s.metrics.DecOrdersBalance()

	return &emptypb.Empty{}, nil
}

// GiveOrderToCustomer обрабатывает выдачу заказов одному клиенту
func (s *Server) GiveOrderToCustomer(ctx context.Context, req *pb.OrderRequestForCustomer) (*emptypb.Empty, error) {
	spanCtx, span := s.tracer.Start(ctx, "GiveOrderToCustomer")
	defer span.End()

	var orderReq orderRequestForCustomerModel
	orderReq.mapFromProto(req)

	if err := orderReq.validate(); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "orderRequestForCustomerModel.validate: %v", err)
	}
	if err := s.service.GiveOrderToCustomer(spanCtx, orderReq.identifier.pvzID, orderReq.customerID, orderReq.identifier.orderID); err != nil {
		if errors.Is(err, orders_errors.ErrNotFound) {
			return &emptypb.Empty{}, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "GiveOrderToCustomer: %v", err)
	}

	s.metrics.IncGivenOrders()
	s.metrics.DecOrdersBalance()

	return &emptypb.Empty{}, nil
}

// TakeRefundFromCustomer обрабатывает возврат заказа клиентом
func (s *Server) TakeRefundFromCustomer(ctx context.Context, req *pb.OrderRequestForCustomer) (*emptypb.Empty, error) {
	spanCtx, span := s.tracer.Start(ctx, "TakeRefundFromCustomer")
	defer span.End()

	var orderReq orderRequestForCustomerModel
	orderReq.mapFromProto(req)

	if err := orderReq.validate(); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "orderRequestForCustomerModel.validate: %v", err)
	}
	if err := s.service.TakeRefundFromCustomer(spanCtx, orderReq.identifier.pvzID, orderReq.customerID, orderReq.identifier.orderID); err != nil {
		if errors.Is(err, orders_errors.ErrNotFound) {
			return &emptypb.Empty{}, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "TakeRefundFromCustomer: %v", err)
	}

	s.metrics.IncRefundedOrders()
	s.metrics.IncOrdersBalance()

	return &emptypb.Empty{}, nil
}

// GetRefundList возвращает страницу возвращенных заказов в этом пвз в виде слайса
// pageNum int - номер страницы, pageSize int - размер страницы
func (s *Server) GetRefundList(ctx context.Context, req *pb.RefundListRequest) (*pb.OrderList, error) {
	spanCtx, span := s.tracer.Start(ctx, "GetRefundList")
	defer span.End()

	var refundListReq refundListRequestModel
	refundListReq.mapFromProto(req)

	if err := refundListReq.validate(); err != nil {
		return &pb.OrderList{}, status.Errorf(codes.InvalidArgument, "refundListRequestModel.validate: %v", err)
	}
	orders, err := s.service.GetRefundList(spanCtx, refundListReq.pvzID, refundListReq.pageNum, refundListReq.pageSize)
	if err != nil {
		if errors.Is(err, orders_errors.ErrNotFound) {
			return &pb.OrderList{}, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return &pb.OrderList{}, status.Errorf(codes.Internal, "GetRefundList: %v", err)
	}
	var protoOrders orderListModel
	protoOrders.mapFromDTO(orders)

	return protoOrders.mapToProto(), nil
}

// GetCustomerOrderList возвращает слайс заказов по ID клиента в этом ПВЗ,
// limit int устанавливает максимальное количество возвращаемых заказов,
// если limit = 0, то ограничения нет
// isInStock bool устанавливает необходимость проверки наличия заказ в пункте,
// в том числе возвращенные
func (s *Server) GetCustomerOrderList(ctx context.Context, req *pb.CustomerOrderListRequest) (*pb.OrderList, error) {
	spanCtx, span := s.tracer.Start(ctx, "GetCustomerOrderList")
	defer span.End()

	var orderListReq customerOrderListRequestModel
	orderListReq.mapFromProto(req)

	if err := orderListReq.validate(); err != nil {
		return &pb.OrderList{}, status.Errorf(codes.InvalidArgument, "customerOrderListRequestModel.validate: %v", err)
	}
	orders, err := s.service.GetCustomerOrderList(spanCtx, orderListReq.pvzID, orderListReq.customerID, orderListReq.limit, orderListReq.isInStock)
	if err != nil {
		if errors.Is(err, orders_errors.ErrNotFound) {
			return &pb.OrderList{}, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return &pb.OrderList{}, status.Errorf(codes.Internal, "GetCustomerOrderList: %v", err)
	}
	var protoOrders orderListModel
	protoOrders.mapFromDTO(orders)

	return protoOrders.mapToProto(), nil
}
