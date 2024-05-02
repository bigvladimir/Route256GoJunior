package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	pvz_errors "homework/internal/app/pvz/errors"
)

// GetByID is a handler for GET method, call GetPvzByID function
func (s *Server) GetByID(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamID]
	keyInt, err := checkKey(key, ok)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pvzJSON, status, err := s.getByIDexec(req.Context(), keyInt)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	w.WriteHeader(status)
	if _, err = w.Write(pvzJSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getByIDexec(ctx context.Context, key int64) ([]byte, int, error) {
	pvzInfo, err := s.service.GetPvzByID(ctx, key)
	if err != nil {
		if errors.Is(err, pvz_errors.ErrNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	pvzJSON, _ := json.Marshal(pvzInfo)

	return pvzJSON, http.StatusOK, nil
}
