package dto

type Pvz struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

type PvzInput struct {
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}
