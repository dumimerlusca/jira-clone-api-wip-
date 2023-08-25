package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Queries struct {
	Db *sql.DB
}

func NewQueries(db *sql.DB) *Queries {
	return &Queries{Db: db}
}

func Connect() (*sql.DB, error) {
	url := os.Getenv("DATABASE_URL")
	fmt.Println("URL", url)
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}
