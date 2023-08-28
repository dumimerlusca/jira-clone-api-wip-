package app

import (
	"jira-clone/packages/db"
	"jira-clone/packages/queries"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var tApp *application
var tHandler http.Handler
var tu *TestUtils

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
	tHandler = tApp.routes()
	tu = &TestUtils{app: tApp}

	os.Exit(m.Run())
}
