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

type user struct {
	Id       *string `json:"id"`
	Username *string `json:"username"`
}

type ticket struct {
	Id           string  `json:"id"`
	Priority     int     `json:"priority"`
	Title        string  `json:"title"`
	Story_points int     `json:"story_points"`
	Description  *string `json:"description"`
	Status       string  `json:"status"`
	Component_id *string `json:"component_id"`
	Created_at   string  `json:"created_at"`
	Updated_at   string  `json:"updated_at"`
	Created_by   *user   `json:"created_by"`
	Assignee     *user   `json:"assignee"`
}

func (app *application) getProjectTickets(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["projectId"]

	rows, err := app.db.Query(`SELECT u.username,u.id, u1.username, u1.id, t.id, t.priority, t.title, t.story_points, t.description,t.status, t.component_id, t.created_at, t.updated_at FROM tickets AS t
	LEFT JOIN users AS u ON u.id=t.assignee_id
	INNER JOIN users AS u1 ON u1.id=t.created_by_id
	WHERE project_id=$1`, projectId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	tickets := []ticket{}

	for rows.Next() {
		var createdBy user
		var assignee user

		t := ticket{Created_by: &createdBy, Assignee: &assignee}

		err := rows.Scan(&assignee.Username, &assignee.Id, &createdBy.Username, &createdBy.Id, &t.Id, &t.Priority, &t.Title, &t.Story_points, &t.Description, &t.Status, &t.Component_id, &t.Created_at, &t.Updated_at)

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
