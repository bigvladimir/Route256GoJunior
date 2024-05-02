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
	pvz_storage "homework/internal/app/pvz/repository/postgresql"
	"homework/internal/app/server"
	"homework/internal/app/server/mdlware"
	"homework/internal/pkg/db"
	"homework/internal/pkg/kafka"
	"homework/internal/pkg/kafkalogger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx, "./docker-compose.yaml")
	if err != nil {
		log.Fatal("db.NewDb:", err)
	}
	defer database.GetPool(ctx).Close()

	brokers, err := kafka.GetBrokers("./config/brokerports.yaml")
	if err != nil {
		log.Fatal("kafka.GetBrokers:", err)
	}

	kafkaConsumer, err := kafka.NewConsumer(brokers)
	if err != nil {
		log.Fatal("kafka.NewConsumer:", err)
	}
	defer func() {
		if err = kafkaConsumer.Close(); err != nil {
			log.Println("kafkaConsumer.Close:", err)
		}
	}()

	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatal("kafka.NewProducer:", err)
	}
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Println("kafkaProducer.Close:", err)
		}
	}()

	validTopicsSet := map[string]string{
		"logs": "",
	}
	// горутина просматривает логи кафки, можно запустить отдельно
	logWatcher := kafkalogger.NewLogWatcher(kafkaConsumer, validTopicsSet)
	go func() {
		err := logWatcher.Subscribe(ctx, "logs")
		if err != nil {
			log.Fatal(`logWatcher.Subscribe(ctx, "logs"):`, err)
		}
	}()

	// шаблон, пока не используется
	orderService := &orders.Service{}
	pvzService := pvz.NewService(pvz_storage.NewPvzStorage(database))
	logger := kafkalogger.NewKafkaLogger(kafkaProducer, "logs")
	service := core.NewCoreService(orderService, pvzService, logger)
	serv := server.NewServer(service)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	checkServerErr := make(chan struct{}, 1)

	http.Handle("/", mdlware.Logger(service, mdlware.Auth(server.CreateRouter(serv))))

	tlsServer := &http.Server{
		Addr:           ":9001",
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// если будет не настроен сертификат и будет ошибка, то программа не закроется и можно будет обращаться по http
	go func() {
		err := tlsServer.ListenAndServeTLS("./config/tls/server.crt", "./config/tls/server.key")
		if err != http.ErrServerClosed {
			log.Println("http.ListenAndServeTLS:", err, "\n Only http requests can be used in this session")
		}
	}()
	defer func() {
		timeoutCtx, timeoutCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCtxCancel()
		if err := tlsServer.Shutdown(timeoutCtx); err != nil {
			log.Println("tlsServer.Shutdown(ctx)", err)
		}
	}()

	httpServer := &http.Server{
		Addr:           ":9000",
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("http.ListenAndServe:", err)
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
