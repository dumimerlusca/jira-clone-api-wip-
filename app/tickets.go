package app

import (
	"encoding/json"
	"jira-clone/packages/events"
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

	ticket, _ := app.queries.FindTicketById(ticketId)

	body.Updated_by_id = userId

	updatedTicket, err := app.queries.UpdateTicket(ticketId, body)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	go func() {
		app.registerTicketUpdatedEvents(ticket, updatedTicket)
	}()

	response.NewSuccessResponse(w, http.StatusOK, updatedTicket)
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

	tickets := []queries.TicketDetails{}

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

		tickets = append(tickets, t)
	}

	response.NewSuccessResponse(w, http.StatusOK, tickets)
}

func (app *application) getTicketDetailsHandler(w http.ResponseWriter, r *http.Request) {
	userId := extractUserId(r)
	ticketKey := mux.Vars(r)["ticketKey"]

	details, err := app.queries.GetTicketDetailsByKeyForUser(ticketKey, userId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, details)
}

type TicketHistoryItem struct {
	events.TicketUpdatedEventData
	Created_at string `json:"created_at"`
}

func (app *application) getTicketHistory(w http.ResponseWriter, r *http.Request) {
	ticketKey := mux.Vars(r)["ticketKey"]

	if ticketKey == "" {
		app.badRequest(w, "invalid ticket id", nil)
		return
	}

	rows, err := app.db.Query(`SELECT created_at, data FROM events WHERE source_id=$1 AND data ->> 'ticketId'=$2 ORDER BY created_at DESC`, events.SourceIdTicketUpdatedEvent, ticketKey)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	history := []*TicketHistoryItem{}

	for rows.Next() {
		var created_at string
		var jData string

		err := rows.Scan(&created_at, &jData)

		if err != nil {
			app.serverError(w, err.Error(), err)
			return
		}

		var eventData TicketHistoryItem

		err = json.Unmarshal([]byte(jData), &eventData)

		if err != nil {
			continue
		}

		eventData.Created_at = created_at

		history = append(history, &eventData)
	}

	response.NewSuccessResponse(w, http.StatusOK, history)
}
