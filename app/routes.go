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

	mux.HandleFunc("/api/projects", app.authMW(app.getProjectsHandler)).Methods("GET")
	mux.HandleFunc("/api/projects/create", app.authMW(app.createProjectHandler)).Methods("POST")
	mux.HandleFunc("/api/projects/details/{projectId}", app.authMW(app.isProjectOwnerMW(app.getProjectDetails))).Methods("GET")
	mux.HandleFunc("/api/projects/update/{projectId}", app.authMW(app.isProjectOwnerMW(app.updateProject))).Methods("PATCH")

	// Project invitations
	mux.HandleFunc("/api/projects/sendInvite/{projectId}", app.authMW(app.isProjectOwnerMW(app.sendProjectInvitationHandler))).Methods("POST")
	mux.HandleFunc("/api/projects/acceptInvite/{inviteId}", app.authMW(app.validateUpdateProjectInviteMW(app.acceptProjectInviteHandler))).Methods("POST")
	mux.HandleFunc("/api/projects/rejectInvite/{inviteId}", app.authMW(app.validateUpdateProjectInviteMW(app.rejectProjectInviteHandler))).Methods("POST")
	mux.HandleFunc("/api/project-invites/sent", app.authMW(app.getSentProjectInvites)).Methods("GET")
	mux.HandleFunc("/api/project-invites/received", app.authMW(app.getReceivedProjectInvites)).Methods("GET")

	return mux
}
