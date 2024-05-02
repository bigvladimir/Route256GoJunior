//go:generate mockgen -source ./server.go -destination=../server/mocks/server.go -package=mock_server
package server

import (
	"context"
	"errors"
	"log"
	"strconv"

	"homework/internal/app/core"
	"homework/internal/app/pvz/dto"
	"homework/internal/pkg/kafkalogger"
)

type coreOps interface {
	AddPvz(ctx context.Context, input dto.PvzInput) (int64, error)
	GetPvzByID(ctx context.Context, id int64) (dto.Pvz, error)
	UpdatePvz(ctx context.Context, input dto.Pvz) error
	DeletePvz(ctx context.Context, id int64) error
	ModifyPvz(ctx context.Context, input dto.Pvz) (int64, error)

	LogMessage(message kafkalogger.Message) error
	LogMessages(message []kafkalogger.Message) error
}

type Server struct {
	service coreOps
}

func NewServer(service *core.Service) *Server {
	return &Server{service: service}
}

func checkKey(key string, ok bool) (int64, error) {
	if !ok {
		err := errors.New("Неправильный ключ запроса")
		log.Println(err)
		return 0, err
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		err := errors.New("Значение ключа не число")
		log.Println(err)
		return 0, err
	}
	if keyInt <= 0 {
		err := errors.New("Неположительный id")
		log.Println(err)
		return 0, err
	}
	return keyInt, nil
}

func validatePvzReq(pvz pvzRequest) error {
	if pvz.Name == "" || pvz.Adress == "" || pvz.Contacts == "" {
		return errors.New("Пустые поля")
	}
	return nil
}

func validateFullPvzReq(pvz pvzFullRequest) error {
	if pvz.ID <= 0 {
		return errors.New("Неположительный id")
	}
	if pvz.Name == "" || pvz.Adress == "" || pvz.Contacts == "" {
		return errors.New("Пустые поля")
	}
	return nil
}
