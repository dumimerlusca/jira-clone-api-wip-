package app

import (
	"database/sql"
	"fmt"
	"jira-clone/packages/db"
	"jira-clone/packages/queries"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
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

	c := cors.New(cors.Options{AllowedOrigins: []string{"*"}, AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"}, AllowedHeaders: []string{"Content-type", "Authorization"}})

	handler := c.Handler(mux)

	log.Fatal(http.ListenAndServe(":3001", handler))
}
