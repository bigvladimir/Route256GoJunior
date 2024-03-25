package repository

type PvzDTO struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Adress   string `db:"adress"`
	Contacts string `db:"contacts"`
}
