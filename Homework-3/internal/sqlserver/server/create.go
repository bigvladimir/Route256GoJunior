package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *Server) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		err = fmt.Errorf("Ошибка чтения тела запроса: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var unm pvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		http.Error(w, "Не удалось десериализовать полученные данные.", http.StatusBadRequest)
		return
	}

	newPvz := ReqToDTO(&unm)
	if pvzIsEmpty(newPvz) {
		http.Error(w, "Пустые поля.", http.StatusBadRequest)
		return
	}

	id, err := s.repo.Add(req.Context(), newPvz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DTOToResp(newPvz, id)
	pvzJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Ошибка при преобразовании данных в json", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(pvzJson)
}
