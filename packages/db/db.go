package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type Queries struct {
	Db *sql.DB
}

func Connect() (*sql.DB, error) {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}
