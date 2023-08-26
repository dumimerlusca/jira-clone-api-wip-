package queries

import "database/sql"

type Queries struct {
	Db *sql.DB
}

func NewQueries(db *sql.DB) *Queries {
	return &Queries{Db: db}
}
