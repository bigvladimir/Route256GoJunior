package grpcserver

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"homework/internal/app/metrics"
	orders_dto "homework/internal/app/orders/dto"
	pvz_dto "homework/internal/app/pvz/dto"
	"homework/internal/pkg/kafkalogger"
	"homework/internal/pkg/pb"
)

type coreOps interface {
	TakeOrderFromCourier(ctx context.Context, order orders_dto.OrderInput) error
	ReturnOrderToCourier(ctx context.Context, pvzID, orderID int) error
	GiveOrderToCustomer(ctx context.Context, pvzID, customerID, orderID int) error
	TakeRefundFromCustomer(ctx context.Context, pvzID, customerID, orderID int) error
	GetRefundList(ctx context.Context, pvzID, pageNum, pageSize int) ([]orders_dto.Order, error)
	GetCustomerOrderList(ctx context.Context, pvzID, customerID, limit int, isInStock bool) ([]orders_dto.Order, error)

	AddPvz(ctx context.Context, input pvz_dto.PvzInput) (int64, error)
	GetPvzByID(ctx context.Context, id int64) (pvz_dto.Pvz, error)
	UpdatePvz(ctx context.Context, input pvz_dto.Pvz) error
	DeletePvz(ctx context.Context, id int64) error
	ModifyPvz(ctx context.Context, input pvz_dto.Pvz) (int64, error)

	LogMessage(message kafkalogger.Message) error
	LogMessages(message []kafkalogger.Message) error
}

type metricsOps interface {
	IncOrdersBalance()
	DecOrdersBalance()
	IncGivenOrders()
	IncRefundedOrders()
	IncReturnedOrders()
}

// TODO возможно стоит добавить интерфейс для трейсера

// Server provides gRPC methods
type Server struct {
	service coreOps
	metrics metricsOps
	tracer  trace.Tracer
	pb.UnimplementedPvzManagerServer
}

// NewServer creates Server
func NewServer(service coreOps, metrics *metrics.Metrics, tracer trace.Tracer) *Server {
	return &Server{
		service: service,
		metrics: metrics,
		tracer:  tracer,
	}
}
