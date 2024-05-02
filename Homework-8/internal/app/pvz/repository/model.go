package repository

import (
	"homework/internal/app/pvz/dto"
)

// PvzRow is the package internal representation of the point of issue of orders
type PvzRow struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Adress   string `db:"adress"`
	Contacts string `db:"contacts"`
}

// MapToDTO converts PvzRow to external representation
func (t *PvzRow) MapToDTO() dto.Pvz {
	return dto.Pvz{
		ID:       t.ID,
		Name:     t.Name,
		Adress:   t.Adress,
		Contacts: t.Contacts,
	}
}

// MapFromDTO converts external representation to PvzRow
func (t *PvzRow) MapFromDTO(d dto.Pvz) {
	t.ID = d.ID
	t.Name = d.Name
	t.Adress = d.Adress
	t.Contacts = d.Contacts
}

// PvzInputRow is the package internal short representation of the point of issue of orders, has no ID
type PvzInputRow struct {
	Name     string `db:"name"`
	Adress   string `db:"adress"`
	Contacts string `db:"contacts"`
}

// MapFromDTO converts external representation to PvzInputRow
func (t *PvzInputRow) MapFromDTO(d dto.PvzInput) {
	t.Name = d.Name
	t.Adress = d.Adress
	t.Contacts = d.Contacts
}
