package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	// TICKETS
	mux.HandleFunc("/api/projects/{projectId}/tickets/create", app.authMW(app.isProjectMemberMW(app.createTicketHandler))).Methods("POST")
	mux.HandleFunc("/api/projects/{projectId}/tickets/update/{ticketId}", app.authMW(app.isProjectMemberMW(app.updateTicketHandler))).Methods("PATCH")

	// AUTH
	mux.HandleFunc("/api/auth/register", app.registerHandler).Methods("POST")
	mux.HandleFunc("/api/auth/login", app.loginHandler).Methods("POST")
	mux.HandleFunc("/api/auth/me", app.authMW(app.getCurrentUserHandler)).Methods("GET")

	// PROJECTS
	mux.HandleFunc("/api/projects", app.authMW(app.getProjectsHandler)).Methods("GET")
	mux.HandleFunc("/api/projects/create", app.authMW(app.createProjectHandler)).Methods("POST")
	mux.HandleFunc("/api/projects/details/{projectId}", app.authMW(app.getProjectDetails)).Methods("GET")
	mux.HandleFunc("/api/projects/update/{projectId}", app.authMW(app.isProjectOwnerMW(app.updateProject))).Methods("PATCH")
	mux.HandleFunc("/api/projects/users/{projectId}", app.authMW(app.getProjectMembers)).Methods("GET")

	// USERS
	mux.HandleFunc("/api/users/details/{username}", app.authMW(app.getUserDetailsHandler)).Methods("GET")

	// Project invitations
	mux.HandleFunc("/api/projects/sendInvite/{projectId}", app.authMW(app.isProjectOwnerMW(app.sendProjectInvitationHandler))).Methods("POST")
	mux.HandleFunc("/api/projects/acceptInvite/{inviteId}", app.authMW(app.validateUpdateProjectInviteMW(app.acceptProjectInviteHandler))).Methods("POST")
	mux.HandleFunc("/api/projects/rejectInvite/{inviteId}", app.authMW(app.validateUpdateProjectInviteMW(app.rejectProjectInviteHandler))).Methods("POST")
	mux.HandleFunc("/api/project-invites/sent", app.authMW(app.getSentProjectInvites)).Methods("GET")
	mux.HandleFunc("/api/project-invites/received", app.authMW(app.getReceivedProjectInvites)).Methods("GET")

	return mux
}
