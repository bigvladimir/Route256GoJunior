package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"

	"homework/internal/app/core"
	"homework/internal/app/grpcserver"
	"homework/internal/app/metrics"
	"homework/internal/app/orders"
	orders_storage "homework/internal/app/orders/repository/postgresql"
	"homework/internal/app/pvz"
	"homework/internal/app/pvz/repository/in_memory_cache"
	pvz_storage "homework/internal/app/pvz/repository/postgresql"
	"homework/internal/pkg/cacheupdater"
	"homework/internal/pkg/config"
	"homework/internal/pkg/db"
	"homework/internal/pkg/db/transaction_manager"
	"homework/internal/pkg/kafka"
	"homework/internal/pkg/kafkalogger"
	"homework/internal/pkg/pb"
	"homework/internal/pkg/tracer"
)

func main() {
	// general initialization
	if err := config.Init(); err != nil {
		log.Fatal("config.Cfg.Init:", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// db initialization
	database, err := db.NewDb(ctx, config.Cfg().GetBdDSN())
	if err != nil {
		log.Fatal("db.NewDb:", err)
	}
	defer database.GetPool(ctx).Close()
	txManager := transaction_manager.NewTransactionManager(database)

	// kafka initialization
	kafkaConsumer, err := kafka.NewConsumer(config.Cfg().GetBrokerPorts())
	if err != nil {
		log.Fatal("kafka.NewConsumer:", err)
	}
	defer func() {
		if err = kafkaConsumer.Close(); err != nil {
			log.Println("kafkaConsumer.Close:", err)
		}
	}()
	kafkaProducer, err := kafka.NewProducer(config.Cfg().GetBrokerPorts())
	if err != nil {
		log.Fatal("kafka.NewProducer:", err)
	}
	defer func() {
		if err = kafkaProducer.Close(); err != nil {
			log.Println("kafkaProducer.Close:", err)
		}
	}()

	// logger initialization
	logger := &kafkalogger.KafkaLogger{}

	// orders storage initialization
	orderStor := orders_storage.NewOrderStorage(txManager)

	// order packages initialization
	packVariants := map[string]orders.PackageVariant{
		orders.BagVariantName:  orders.NewBagPackage(),
		orders.BoxVariantName:  orders.NewBoxPackage(),
		orders.FilmVariantName: orders.NewFilmPackage(),
	}

	// pvz storage initialization
	validCacheTopicsSet := map[string]string{
		"pvz_cache_updater": "",
	}
	cacheUpdateWriter := cacheupdater.NewCacheUpdateWriter(kafkaProducer, "pvz_cache_updater")
	cacheUpdateReader := cacheupdater.NewCacheUpdateReader(kafkaConsumer, validCacheTopicsSet)
	if err != nil {
		log.Fatal("in_memory_cache.NewInMemoryCache", err)
	}
	pvzStor := pvz_storage.NewPvzStorage(txManager)
	inMemoryCache, err := in_memory_cache.NewInMemoryCache(ctx, pvzStor, cacheUpdateReader)

	// service initialization
	orderService := orders.NewService(orderStor, packVariants)
	pvzService := pvz.NewService(pvzStor, txManager, inMemoryCache, cacheUpdateWriter)
	service := core.NewCoreService(orderService, pvzService, logger)

	// metrics initialization
	serverMetrics, reg := metrics.InitMetrics()

	// tracer initialization
	tracerShutdown, err := tracer.InitProvider("pvz-manager", "localhost:16686")
	if err != nil {
		log.Fatal("tracer.InitProvider", err)
	}
	defer func() {
		if tracerShutdownErr := tracerShutdown(ctx); tracerShutdownErr != nil {
			log.Println("tracerShutdown:", tracerShutdownErr)
		}
	}()

	// server initialization
	serverImplementation := grpcserver.NewServer(service, serverMetrics, otel.Tracer("pvz-manager"))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9400))
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(serverMetrics.StandartGrpcMetrics.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(serverMetrics.StandartGrpcMetrics.StreamServerInterceptor()),
	)

	pb.RegisterPvzManagerServer(grpcServer, serverImplementation)
	serverMetrics.StandartGrpcMetrics.InitializeMetrics(grpcServer)

	serverForMetrics := &http.Server{
		Addr:           ":9401",
		Handler:        promhttp.HandlerFor(reg, promhttp.HandlerOpts{EnableOpenMetrics: true}),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if serverForMetricsErr := serverForMetrics.ListenAndServe(); serverForMetricsErr != http.ErrServerClosed {
			log.Fatal("serverForMetrics.ListenAndServe:", serverForMetricsErr)
		}
	}()
	defer func() {
		timeoutCtx, timeoutCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCtxCancel()
		if serverForMetricsErr := serverForMetrics.Shutdown(timeoutCtx); serverForMetricsErr != nil {
			log.Println("httpServer.Shutdown(ctx)", serverForMetricsErr)
		}
	}()

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal("grpcServer.Serve:", err)
	}
}
