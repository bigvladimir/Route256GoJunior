package server

import (
	"errors"
	"homework/internal/sqlserver/repository"
	"net/http"
)

func (s *Server) Delete(w http.ResponseWriter, req *http.Request) {
	keyInt, err := checkKey(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.repo.Delete(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
