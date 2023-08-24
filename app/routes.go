package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/api/tickets", app.createTicketHandler).Methods("POST")
	mux.HandleFunc("/api/tickets/{id}", app.updateTicketHandler).Methods("PATCH")

	mux.HandleFunc("/api/auth/register", app.registerHandler).Methods("POST")
	mux.HandleFunc("/api/auth/login", app.loginHandler).Methods("POST")

	mux.HandleFunc("/api/projects", app.createProjectHandler).Methods("POST")
	mux.HandleFunc("/api/projects", app.getProjectsHandler).Methods("GET")
	mux.HandleFunc("/api/projects/{projectId}", app.getProjectDetails).Methods("GET")
	mux.HandleFunc("/api/projects/{projectId}", app.authMW(app.updateProject)).Methods("PATCH")
	return mux
}
