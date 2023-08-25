package db

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env.test")

	if err != nil {
		log.Fatal("failed to load env variables", err)
		return
	}

	db, err := Connect()

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQueries = NewQueries(db)

	os.Exit(m.Run())
}
