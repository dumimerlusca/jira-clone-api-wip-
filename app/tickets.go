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
	userId := r.Context().Value(ContextKey("userId")).(string)

	var body queries.UpdateTicketDTO

	err := util.ReadAndUnmarshal(r.Body, &body)

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	body.Updated_by_id = &userId

	ticket, err := app.queries.UpdateTicket(ticketId, body)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, ticket)
}

func (app *application) getProjectTickets(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["projectId"]

	rows, err := app.db.Query(`SELECT 
		id,
		key,
		type,
		priority,
		title,
		story_points,
		description,
		status,
		component_id,
		created_at,
		updated_at,
		creator_id,
		creator_username,
		assignee_id,
		assignee_username
	FROM tickets_view WHERE project_id=$1`, projectId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	tickets := []TicketItem{}

	for rows.Next() {
		var createdBy UserItem
		var assignee UserItem

		t := TicketItem{Creator: &createdBy, Assignee: &assignee}

		err := rows.Scan(&t.Id, &t.Key, &t.Type, &t.Priority, &t.Title, &t.Story_points, &t.Description, &t.Status, &t.Component_id, &t.Created_at, &t.Updated_at, &t.Creator.Id, &t.Creator.Username, &t.Assignee.Id, &t.Assignee.Username)

		if err != nil {
			app.serverError(w, err.Error(), err)
			return
		}

		if t.Assignee.Id == nil {
			t.Assignee = nil
		}

		tickets = append(tickets, t)
	}

	response.NewSuccessResponse(w, http.StatusOK, tickets)
}
