package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Modify is a handler for PUT method, call ModifyPvz function
func (s *Server) Modify(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		err = fmt.Errorf("Ошибка чтения тела запроса: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var unm pvzFullRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		http.Error(w, "Не удалось десериализовать полученные данные", http.StatusBadRequest)
		return
	}

	err = validateFullPvzReq(unm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pvzJSON, status, err := s.modifyExec(req.Context(), unm)
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

func (s *Server) modifyExec(ctx context.Context, req pvzFullRequest) ([]byte, int, error) {
	newPvz := req.mapToPvz()

	id, err := s.service.ModifyPvz(ctx, newPvz)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var status int
	if id <= 0 {
		id = newPvz.ID
		status = http.StatusOK
	} else {
		status = http.StatusCreated
	}
	var resp pvzResponse
	resp.mapFromPvz(newPvz)
	resp.ID = id
	pvzJSON, _ := json.Marshal(resp)

	return pvzJSON, status, nil
}
