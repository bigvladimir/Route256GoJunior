package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *Server) Update(w http.ResponseWriter, req *http.Request) {
	keyInt, err := checkKey(req)
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
		http.Error(w, "Не удалось десериализовать полученные данные.", http.StatusBadRequest)
		return
	}

	newPvz := ReqToDTO(&unm)
	if pvzIsEmpty(newPvz) {
		http.Error(w, "Пустые поля.", http.StatusBadRequest)
		return
	}

	newPvz.ID = keyInt

	err = s.repo.Update(req.Context(), newPvz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
