package server

import (
	"homework/internal/sqlserver/repository"
)

type pvzRequest struct {
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

type pvzFullRequest struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

type pvzResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

func ReqToDTO(req *pvzRequest) *repository.PvzDTO {
	return &repository.PvzDTO{
		Name:     req.Name,
		Adress:   req.Adress,
		Contacts: req.Contacts,
	}
}

func DTOToResp(req *repository.PvzDTO, id int64) *pvzResponse {
	return &pvzResponse{
		ID:       id,
		Name:     req.Name,
		Adress:   req.Adress,
		Contacts: req.Contacts,
	}
}

func FullReqToDTO(req *pvzFullRequest) *repository.PvzDTO {
	return &repository.PvzDTO{
		ID:       req.ID,
		Name:     req.Name,
		Adress:   req.Adress,
		Contacts: req.Contacts,
	}
}
