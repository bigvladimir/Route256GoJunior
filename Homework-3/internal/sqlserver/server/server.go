package server

import (
	"errors"
	"homework/internal/sqlserver/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Server struct {
	repo *repository.PvzStore
}

func NewServer(setrepo *repository.PvzStore) Server {
	return Server{repo: setrepo}
}

func pvzIsEmpty(pvz *repository.PvzDTO) bool {
	if pvz.Name == "" || pvz.Adress == "" || pvz.Contacts == "" {
		return true
	}
	return false
}

func checkKey(req *http.Request) (int64, error) {
	key, ok := mux.Vars(req)[queryParamID]
	if !ok {
		err := errors.New("Неправильный ключ запроса.")
		log.Println(err)
		return -1, err
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		err := errors.New("Значение ключа не число.")
		log.Println(err)
		return -1, err
	}
	if keyInt <= 0 {
		err := errors.New("Неположительный id.")
		log.Println(err)
		return -1, err
	}
	return keyInt, nil
}
