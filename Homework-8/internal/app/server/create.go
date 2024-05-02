package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Create is a handler for POST method, call AddPvz function
func (s *Server) Create(w http.ResponseWriter, req *http.Request) {
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

	pvzJSON, status, err := s.createExec(req.Context(), unm)
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

func (s *Server) createExec(ctx context.Context, req pvzRequest) ([]byte, int, error) {
	newPvz := req.mapToPvzInput()

	id, err := s.service.AddPvz(ctx, newPvz)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var resp pvzResponse
	resp.mapFromPvzInput(newPvz)
	resp.ID = id
	pvzJSON, _ := json.Marshal(resp)

	return pvzJSON, http.StatusCreated, nil
}
