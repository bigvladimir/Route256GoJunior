package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"homework/internal/app/core"
	"homework/internal/app/orders"
	"homework/internal/app/pvz"
	"homework/internal/app/pvz/repository/in_memory_cache"
	pvz_storage "homework/internal/app/pvz/repository/postgresql"
	"homework/internal/app/server"
	"homework/internal/app/server/mdlware"
	"homework/internal/pkg/cacheupdater"
	"homework/internal/pkg/config"
	"homework/internal/pkg/db"
	"homework/internal/pkg/db/transaction_manager"
	"homework/internal/pkg/kafka"
	"homework/internal/pkg/kafkalogger"
)

func main() {
	// general initialization
	if err := config.Init(); err != nil {
		log.Fatal("config.Cfg.Init:", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
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
	validLogTopicsSet := map[string]string{
		"logs": "",
	}
	logger := kafkalogger.NewKafkaLogger(kafkaProducer, "logs")
	logWatcher := kafkalogger.NewLogWatcher(kafkaConsumer, validLogTopicsSet)
	if err = logWatcher.Subscribe(ctx, "logs"); err != nil {
		log.Fatal(`logWatcher.Subscribe(ctx, "logs"):`, err)
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
	orderService := &orders.Service{}
	pvzService := pvz.NewService(pvzStor, txManager, inMemoryCache, cacheUpdateWriter)
	service := core.NewCoreService(orderService, pvzService, logger)

	// signals initialization
	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	checkServerErr := make(chan struct{}, 1)
	defer close(checkServerErr)

	// server initialization
	serv := server.NewServer(service)
	http.Handle("/", mdlware.Logger(service, mdlware.Auth(server.CreateRouter(serv))))
	httpServer := &http.Server{
		Addr:           ":9000",
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("http.ListenAndServe:", err)
			checkServerErr <- struct{}{}
		}
	}()
	defer func() {
		timeoutCtx, timeoutCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCtxCancel()
		if err := httpServer.Shutdown(timeoutCtx); err != nil {
			log.Println("httpServer.Shutdown(ctx)", err)
		}
	}()
	log.Println("HTTP server started")

	// exit loop
Exit:
	for {
		select {
		case <-checkServerErr:
			break Exit
		case <-sig:
			break Exit
		}
	}
}
