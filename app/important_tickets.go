package app

import (
	"jira-clone/packages/queries"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) getMyImportantTicketsHandler(w http.ResponseWriter, r *http.Request) {
	userId := extractUserId(r)

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
		FROM important_tickets
		INNER JOIN tickets_view ON tickets_view.id = ticket_id
		WHERE user_id=$1`, userId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	tickets := []*queries.TicketDetails{}

	for rows.Next() {
		var createdBy queries.UserItem
		var assignee queries.UserItem

		t := queries.TicketDetails{Creator: &createdBy, Assignee: &assignee}

		err := rows.Scan(&t.Id, &t.Key, &t.Type, &t.Priority, &t.Title, &t.Story_points, &t.Description, &t.Status, &t.Component_id, &t.Created_at, &t.Updated_at, &t.Creator.Id, &t.Creator.Username, &t.Assignee.Id, &t.Assignee.Username)

		if err != nil {
			app.serverError(w, err.Error(), err)
			return
		}

		if t.Assignee.Id == nil {
			t.Assignee = nil
		}

		tickets = append(tickets, &t)
	}

	response.NewSuccessResponse(w, http.StatusOK, tickets)
}

type AddImportantTicketPayload struct {
	TicketId string `json:"ticketId"`
}

func (app *application) addImportantTicketHandler(w http.ResponseWriter, r *http.Request) {
	userId := extractUserId(r)

	var body AddImportantTicketPayload

	err := util.ReadAndUnmarshal(r.Body, &body)

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	ticket, err := app.queries.FindTicketById(body.TicketId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	_, err = app.queries.CreateImportantTicket(userId, body.TicketId, ticket.Project_id)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusCreated, struct{}{})
}

func (app *application) deleteImportantTicketHandler(w http.ResponseWriter, r *http.Request) {
	userId := extractUserId(r)
	ticketId := mux.Vars(r)["ticketId"]

	err := app.queries.DeleteImportantTicket(userId, ticketId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, nil)
}
