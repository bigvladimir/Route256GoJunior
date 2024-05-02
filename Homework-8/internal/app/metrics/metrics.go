package metrics

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics contains all server metrics
type Metrics struct {
	StandartGrpcMetrics *grpc_prometheus.ServerMetrics

	ordersOnPvzBalance      prometheus.Gauge
	givenOrders             prometheus.Counter
	refundedOrders          prometheus.Counter
	ordersReturnedToCourier prometheus.Counter
}

// InitMetrics initialises and register all server metrics
func InitMetrics() (*Metrics, *prometheus.Registry) {
	reg := prometheus.NewRegistry()

	metrics := &Metrics{
		StandartGrpcMetrics: grpc_prometheus.NewServerMetrics(),
		ordersOnPvzBalance: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "orders_on_PVZ_balance",
			Help: "Total number of orders in all PVZ.",
		}),
		givenOrders: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "given_orders",
			Help: "Total number of given orders.",
		}),
		refundedOrders: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "redunded_orders",
			Help: "Total number of refunded orders.",
		}),
		ordersReturnedToCourier: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_returned_to_courier",
			Help: "Total number of orders returned to courier.",
		}),
	}

	reg.MustRegister(
		metrics.StandartGrpcMetrics,

		metrics.givenOrders,
		metrics.ordersOnPvzBalance,
		metrics.refundedOrders,
		metrics.ordersReturnedToCourier,
	)

	return metrics, reg
}

// IncOrdersBalance increase by 1 ordersOnPvzBalance metric
func (m *Metrics) IncOrdersBalance() {
	m.ordersOnPvzBalance.Inc()
}

// DecOrdersBalance decrease by 1 ordersOnPvzBalance metric
func (m *Metrics) DecOrdersBalance() {
	m.ordersOnPvzBalance.Dec()
}

// IncGivenOrders increase by 1 givenOrders metric
func (m *Metrics) IncGivenOrders() {
	m.givenOrders.Inc()
}

// IncRefundedOrders increase by 1 refundedOrders metric
func (m *Metrics) IncRefundedOrders() {
	m.refundedOrders.Inc()
}

// IncReturnedOrders increase by 1 ordersReturnedToCourier metric
func (m *Metrics) IncReturnedOrders() {
	m.ordersReturnedToCourier.Inc()
}
