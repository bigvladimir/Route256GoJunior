package server

import (
	pvz_dto "homework/internal/app/pvz/dto"
)

type pvzRequest struct {
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

func (t *pvzRequest) mapToPvz() pvz_dto.Pvz {
	return pvz_dto.Pvz{
		Name:     t.Name,
		Adress:   t.Adress,
		Contacts: t.Contacts,
	}
}

func (t *pvzRequest) mapToPvzInput() pvz_dto.PvzInput {
	return pvz_dto.PvzInput{
		Name:     t.Name,
		Adress:   t.Adress,
		Contacts: t.Contacts,
	}
}

type pvzFullRequest struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

func (t *pvzFullRequest) mapToPvz() pvz_dto.Pvz {
	return pvz_dto.Pvz{
		ID:       t.ID,
		Name:     t.Name,
		Adress:   t.Adress,
		Contacts: t.Contacts,
	}
}

type pvzResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

func (t *pvzResponse) mapFromPvz(req pvz_dto.Pvz) {
	t.ID = req.ID
	t.Name = req.Name
	t.Adress = req.Adress
	t.Contacts = req.Contacts
}

func (t *pvzResponse) mapFromPvzInput(req pvz_dto.PvzInput) {
	t.Name = req.Name
	t.Adress = req.Adress
	t.Contacts = req.Contacts
}
