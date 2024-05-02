package repository

import (
	"homework/internal/app/pvz/dto"
)

type PvzRow struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Adress   string `db:"adress"`
	Contacts string `db:"contacts"`
}

func (t *PvzRow) MapToDTO() dto.Pvz {
	return dto.Pvz{
		ID:       t.ID,
		Name:     t.Name,
		Adress:   t.Adress,
		Contacts: t.Contacts,
	}
}

func (t *PvzRow) MapFromDTO(d dto.Pvz) {
	t.ID = d.ID
	t.Name = d.Name
	t.Adress = d.Adress
	t.Contacts = d.Contacts
}

type PvzInputRow struct {
	Name     string `db:"name"`
	Adress   string `db:"adress"`
	Contacts string `db:"contacts"`
}

func (t *PvzInputRow) MapFromDTO(d dto.PvzInput) {
	t.Name = d.Name
	t.Adress = d.Adress
	t.Contacts = d.Contacts
}
