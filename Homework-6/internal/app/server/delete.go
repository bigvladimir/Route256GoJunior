package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	pvz_errors "homework/internal/app/pvz/errors"
)

// Delete is a handler for DELETE method, call DeletePvz function
func (s *Server) Delete(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamID]
	keyInt, err := checkKey(key, ok)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := s.deleteExec(req.Context(), keyInt)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	w.WriteHeader(status)
}

func (s *Server) deleteExec(ctx context.Context, key int64) (int, error) {
	err := s.service.DeletePvz(ctx, key)
	if err != nil {
		if errors.Is(err, pvz_errors.ErrNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
