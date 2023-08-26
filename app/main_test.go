package app

import (
	"jira-clone/packages/db"
	"jira-clone/packages/queries"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var tApp *application

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env.test")

	if err != nil {
		log.Fatal("failed to load env variables", err)
		return
	}

	db, err := db.Connect()

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	queries := queries.NewQueries(db)
	tApp = &application{db: db, queries: queries}

	os.Exit(m.Run())
}
