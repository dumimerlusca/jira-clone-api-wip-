package app

import (
	"database/sql"
	"fmt"
	"jira-clone/packages/db"
	"jira-clone/packages/queries"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

type application struct {
	db      *sql.DB
	queries *queries.Queries
}

func Init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	sqlDb, err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database")

	app := application{db: sqlDb, queries: queries.NewQueries(sqlDb)}

	mux := app.routes()

	http.ListenAndServe(":3001", mux)
}
