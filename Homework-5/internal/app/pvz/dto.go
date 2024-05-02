package pvz

type Pvz struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

func (t *Pvz) mapToModel() pvzRow {
	return pvzRow{
		ID:       t.ID,
		Name:     t.Name,
		Adress:   t.Adress,
		Contacts: t.Contacts,
	}
}

type PvzInput struct {
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

func (t *PvzInput) mapToModel() pvzInputRow {
	return pvzInputRow{
		name:     t.Name,
		adress:   t.Adress,
		contacts: t.Contacts,
	}
}
