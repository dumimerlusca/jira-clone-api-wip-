package queries

import (
	"jira-clone/packages/db"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var tQueries *Queries

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env.test")

	if err != nil {
		log.Fatal("failed to load env variables", err)
		return
	}

	db, err := db.Connect()

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	tQueries = NewQueries(db)

	os.Exit(m.Run())
}
