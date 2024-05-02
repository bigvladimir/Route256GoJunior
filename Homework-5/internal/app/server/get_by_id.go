package server

import (
	"context"
	"encoding/json"
	"errors"
	"homework/internal/app/pvz"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) GetByID(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamID]
	keyInt, err := checkKey(key, ok)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pvzJson, status, err := s.getByIDexec(req.Context(), keyInt)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	w.WriteHeader(status)
	w.Write(pvzJson)
}

func (s *Server) getByIDexec(ctx context.Context, key int64) ([]byte, int, error) {
	pvzInfo, err := s.service.GetPvzByID(ctx, key)
	if err != nil {
		if errors.Is(err, pvz.ErrNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	pvzJson, _ := json.Marshal(pvzInfo)

	return pvzJson, http.StatusOK, nil
}
