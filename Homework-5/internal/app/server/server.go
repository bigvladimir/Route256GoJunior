package server

import (
	"errors"
	"homework/internal/app/core"
	"log"
	"strconv"
)

type Server struct {
	service core.CoreOps
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
