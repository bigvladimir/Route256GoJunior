package dto

// Pvz is a representation of the point of issue of orders
type Pvz struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}

// PvzInput is a short representation of the point of issue of orders, has no ID
type PvzInput struct {
	Name     string `json:"name"`
	Adress   string `json:"adress"`
	Contacts string `json:"contacts"`
}
