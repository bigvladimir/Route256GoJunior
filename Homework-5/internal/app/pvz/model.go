package pvz

type pvzRow struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Adress   string `db:"adress"`
	Contacts string `db:"contacts"`
}

func (t *pvzRow) mapToDTO() Pvz {
	return Pvz{
		ID:       t.ID,
		Name:     t.Name,
		Adress:   t.Adress,
		Contacts: t.Contacts,
	}
}

type pvzInputRow struct {
	name     string `db:"name"`
	adress   string `db:"adress"`
	contacts string `db:"contacts"`
}
