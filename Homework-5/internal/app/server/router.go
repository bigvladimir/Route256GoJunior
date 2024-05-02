package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// создание обработчика всех доступных уровней
func CreateRouter(serv *Server) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/%s", BasePvzPath), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			serv.Create(w, req)
		case http.MethodPut:
			serv.Modify(w, req)
		default:
			http.Error(w, "Метод не обрабатывается сервером", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(fmt.Sprintf("/%s/{%s:[0-9]+}", BasePvzPath, queryParamID), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			serv.GetByID(w, req)
		case http.MethodDelete:
			serv.Delete(w, req)
		case http.MethodPatch:
			serv.Update(w, req)
		default:
			http.Error(w, "Метод не обрабатывается сервером", http.StatusMethodNotAllowed)
		}
	})

	return router
}
