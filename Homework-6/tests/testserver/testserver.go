//go:build integration

package testserver

import (
	"net/http"

	"homework/internal/app/core"
	"homework/internal/app/orders"
	"homework/internal/app/pvz"
	pvz_storage "homework/internal/app/pvz/repository/postgresql"
	"homework/internal/app/server"
	"homework/internal/pkg/kafkalogger"
	"homework/tests/postgresql"
)

type Handlers interface {
	Create(w http.ResponseWriter, req *http.Request)
	Modify(w http.ResponseWriter, req *http.Request)
	GetByID(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
}

type TServer struct {
	H Handlers
}

func NewTServer(db *postgresql.TDB) *TServer {
	orderService := &orders.Service{}
	pvzService := pvz.NewService(pvz_storage.NewPvzStorage(db.DB))
	logger := &kafkalogger.KafkaLogger{}
	serv := server.NewServer(core.NewCoreService(orderService, pvzService, logger))
	return &TServer{H: serv}
}
