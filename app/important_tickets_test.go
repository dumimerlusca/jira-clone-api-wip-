package app

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddImportantTicketHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, "POST", `/api/important-tickets/create`, struct{}{})

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should return 201 status and save the ticket to database", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)

		res, _ := tu.SendAuthorizedReq(t, "POST", "/api/important-tickets/create", AddImportantTicketPayload{TicketId: ticket.Id}, ticket.Created_by_id)

		tu.RequireStatus(t, res, http.StatusCreated)

		item, _ := tu.app.queries.FindImportantTicket(ticket.Created_by_id, ticket.Id)

		assert.NotNil(t, item)
	})
}

func TestGetMyImportantTicketsHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, "GET", `/api/important-tickets`, struct{}{})

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should return 200 and list of tickets", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)

		tu.SendAuthorizedReq(t, "POST", "/api/important-tickets/create", AddImportantTicketPayload{TicketId: ticket.Id}, ticket.Created_by_id)

		res, _ := tu.SendAuthorizedReq(t, "GET", "/api/important-tickets", struct{}{}, ticket.Created_by_id)

		tu.RequireStatus(t, res, http.StatusOK)

		data := tu.GetSuccessResponseData(t, res).([]any)

		assert.Equal(t, 1, len(data))

		item := data[0].(map[string]any)

		assert.Equal(t, ticket.Id, item["id"])
		assert.Equal(t, ticket.Title, item["title"])

	})
}

func TestDeleteImportantTicketHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, "DELETE", `/api/important-tickets/delete/1`, struct{}{})

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should return 200 and properly delete the ticket from the table", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)

		tu.SendAuthorizedReq(t, "POST", "/api/important-tickets/create", AddImportantTicketPayload{TicketId: ticket.Id}, ticket.Created_by_id)

		res, _ := tu.SendAuthorizedReq(t, "DELETE", "/api/important-tickets/delete/"+ticket.Id, struct{}{}, ticket.Created_by_id)

		tu.RequireStatus(t, res, http.StatusOK)

		item, _ := tu.app.queries.FindImportantTicket(ticket.Created_by_id, ticket.Id)

		assert.Nil(t, item)
	})
}
