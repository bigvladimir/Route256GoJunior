package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework/internal/app/pvz"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Update(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamID]
	keyInt, err := checkKey(key, ok)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		err = fmt.Errorf("Ошибка чтения тела запроса: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var unm pvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		http.Error(w, "Не удалось десериализовать полученные данные", http.StatusBadRequest)
		return
	}

	err = validatePvzReq(unm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := s.updateExec(req.Context(), keyInt, unm)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	w.WriteHeader(status)
}

func (s *Server) updateExec(ctx context.Context, key int64, req pvzRequest) (int, error) {
	newPvz := req.mapToPvz()
	newPvz.ID = key

	err := s.service.UpdatePvz(ctx, newPvz)
	if err != nil {
		if errors.Is(err, pvz.ErrNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
