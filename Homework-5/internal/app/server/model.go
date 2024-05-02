package server

import (
	"homework/internal/app/pvz"
)

type pvzRequest struct {
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

func (t *pvzRequest) mapToPvz() pvz.Pvz {
	return pvz.Pvz{
		Name:     t.Name,
		Adress:   t.Adress,
		Contacts: t.Contacts,
	}
}

func (t *pvzRequest) mapToPvzInput() pvz.PvzInput {
	return pvz.PvzInput{
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

func (t *pvzFullRequest) mapToPvz() pvz.Pvz {
	return pvz.Pvz{
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

func (t *pvzResponse) mapFromPvz(req pvz.Pvz) {
	t.ID = req.ID
	t.Name = req.Name
	t.Adress = req.Adress
	t.Contacts = req.Contacts
}

func (t *pvzResponse) mapFromPvzInput(req pvz.PvzInput) {
	t.Name = req.Name
	t.Adress = req.Adress
	t.Contacts = req.Contacts
}
