package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *Server) Modify(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		err = fmt.Errorf("Ошибка чтения тела запроса: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var unm pvzFullRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		http.Error(w, "Не удалось десериализовать полученные данные.", http.StatusBadRequest)
		return
	}

	newPvz := FullReqToDTO(&unm)
	if newPvz.ID <= 0 {
		http.Error(w, "Неположительный id.", http.StatusBadRequest)
		return
	}
	if pvzIsEmpty(newPvz) {
		http.Error(w, "Пустые поля.", http.StatusBadRequest)
		return
	}

	id, err := s.repo.Modify(req.Context(), newPvz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id <= 0 {
		id = newPvz.ID
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	resp := DTOToResp(newPvz, id)
	pvzJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Ошибка при преобразовании данных в json", http.StatusInternalServerError)
		return
	}

	w.Write(pvzJson)
}
