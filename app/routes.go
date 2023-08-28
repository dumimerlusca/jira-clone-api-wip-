package app

import (
	c "jira-clone/packages/consts"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc(c.ApiPathCreateTicker, app.createTicketHandler).Methods("POST")
	mux.HandleFunc(c.ApiPathUpdateTicket, app.updateTicketHandler).Methods("PATCH")

	mux.HandleFunc(c.ApiPathRegister, app.registerHandler).Methods("POST")
	mux.HandleFunc(c.ApiPathLogin, app.loginHandler).Methods("POST")

	mux.HandleFunc("/api/projects", app.authMW(app.getProjectsHandler)).Methods("GET")
	mux.HandleFunc(c.ApiPathCreateProject, app.authMW(app.createProjectHandler)).Methods("POST")
	mux.HandleFunc(c.ApiPathGetProjectDetails, app.authMW(app.isProjectOwnerMW(app.getProjectDetails))).Methods("GET")
	mux.HandleFunc(c.ApiPathUpdateProject, app.authMW(app.isProjectOwnerMW(app.updateProject))).Methods("PATCH")

	// Project invitations
	mux.HandleFunc(c.ApiPathSendProjectInvite, app.authMW(app.isProjectOwnerMW(app.sendProjectInvitationHandler))).Methods("POST")
	mux.HandleFunc(c.ApiPathAcceptProjectInvite, app.authMW(app.validateUpdateProjectInviteMW(app.acceptProjectInviteHandler))).Methods("POST")
	mux.HandleFunc(c.ApiPathRejectProjectInvite, app.authMW(app.validateUpdateProjectInviteMW(app.rejectProjectInviteHandler))).Methods("POST")
	mux.HandleFunc(c.ApiPathGetSentProjectInvites, app.authMW(app.getSentProjectInvites)).Methods("GET")
	mux.HandleFunc(c.ApiPathGetReceivedProjectInvites, app.authMW(app.getReceivedProjectInvites)).Methods("GET")

	return mux
}
