package app

import (
	"jira-clone/packages/queries"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) createTicketHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(ContextKey("userId")).(string)
	projectId := mux.Vars(r)["projectId"]

	var body queries.CreateTicketDTO

	err := util.ReadAndUnmarshal(r.Body, &body)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	err = body.Validate()

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	body.Created_by_id = userId
	body.Project_id = projectId

	ticket, err := app.queries.CreateTicket(body)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusCreated, ticket)

}

func (app *application) updateTicketHandler(w http.ResponseWriter, r *http.Request) {
	ticketId := mux.Vars(r)["ticketId"]

	var body queries.UpdateTicketDTO

	err := util.ReadAndUnmarshal(r.Body, &body)

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	ticket, err := app.queries.UpdateTicket(ticketId, body)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, ticket)
}
