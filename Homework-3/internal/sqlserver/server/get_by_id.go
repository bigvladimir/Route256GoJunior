package server

import (
	"encoding/json"
	"errors"
	"homework/internal/sqlserver/repository"
	"net/http"
)

func (s *Server) GetByID(w http.ResponseWriter, req *http.Request) {
	keyInt, err := checkKey(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pvz, err := s.repo.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pvzJson, err := json.Marshal(pvz)
	if err != nil {
		http.Error(w, "Ошибка при преобразовании данных в json", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(pvzJson)
}
