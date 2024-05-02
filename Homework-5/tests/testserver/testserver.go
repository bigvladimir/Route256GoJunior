//go:build integration

package testserver

import (
	"homework/internal/app/core"
	"homework/internal/app/orders"
	"homework/internal/app/pvz"
	"homework/internal/app/server"
	"homework/tests/postgresql"
	"net/http"
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
	orderService := orders.NewService(orders.NewOrderStorage(db.DB), map[string]orders.PackageVariant{})
	pvzService := pvz.NewService(pvz.NewPvzStorage(db.DB))
	serv := server.NewServer(core.NewCoreService(orderService, pvzService))
	return &TServer{H: serv}
}
