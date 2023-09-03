package app

import (
	"fmt"
	"jira-clone/packages/queries"
	"jira-clone/packages/random"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCommentHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, "POST", "/api/tickets/ticket1/comments/create", nil)

		tu.RequireStatus(t, res, http.StatusUnauthorized)

	})

	t.Run("require user to be project member", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)
		user, _ := tu.app.queries.CreateRandomUser(t)
		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, fmt.Sprintf("/api/tickets/%v/comments/create", ticket.Id), nil, user.Id)
		tu.RequireStatus(t, res, http.StatusUnauthorized)

		res, _ = tu.SendAuthorizedReq(t, http.MethodPost, fmt.Sprintf("/api/tickets/%v/comments/create", ticket.Id), nil, ticket.Created_by_id)

		require.NotEqual(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("should return 201 status and save comment to database", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)

		payload := createCommentRequestPayload{Text: random.RandomString(40)}

		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, fmt.Sprintf("/api/tickets/%v/comments/create", ticket.Id), payload, ticket.Created_by_id)

		tu.RequireStatus(t, res, http.StatusCreated)
		tu.IsSuccessResponseWithData(t, res)
	})
}

func TestGetTicketCommentsHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, "GET", "/api/tickets/ticket1/comments", nil)

		tu.RequireStatus(t, res, http.StatusUnauthorized)

	})

	t.Run("require user to be project member", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)
		user, _ := tu.app.queries.CreateRandomUser(t)
		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, fmt.Sprintf("/api/tickets/%v/comments", ticket.Id), nil, user.Id)
		tu.RequireStatus(t, res, http.StatusUnauthorized)

		res, _ = tu.SendAuthorizedReq(t, http.MethodGet, fmt.Sprintf("/api/tickets/%v/comments", ticket.Id), nil, ticket.Created_by_id)

		require.NotEqual(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("should return 200 status and array of comments to clients", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)
		dto := queries.CreateCommentDTO{Author_id: ticket.Created_by_id, Text: random.RandomString(30), Ticket_id: ticket.Id}

		count := 5

		for i := 0; i < count; i++ {
			tu.app.queries.CreateComment(dto)
		}

		res, _ := tu.SendAuthorizedReq(t, "GET", fmt.Sprintf("/api/tickets/%v/comments", ticket.Id), nil, ticket.Created_by_id)

		tu.RequireStatus(t, res, http.StatusOK)

		var resBody response.SuccessResponse

		util.ReadAndUnmarshal(res.Body, &resBody)
		data, ok := resBody.Data.([]any)

		assert.Equal(t, true, ok)
		assert.Equal(t, count, len(data))

	})
}

func TestDeleteCommentHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, "DELETE", "/api/comments/delete/1", nil)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})
	t.Run("should return 401 status if user is not comment author", func(t *testing.T) {
		com := tu.app.queries.CreateRandomComment(t)
		user, _ := tu.app.queries.CreateRandomUser(t)

		res, _ := tu.SendAuthorizedReq(t, "DELETE", fmt.Sprintf("/api/comments/delete/%v", com.Id), nil, user.Id)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should return status 204 and delete comment from database", func(t *testing.T) {
		com := tu.app.queries.CreateRandomComment(t)

		res, _ := tu.SendAuthorizedReq(t, "DELETE", fmt.Sprintf("/api/comments/delete/%v", com.Id), nil, com.Author_id)

		tu.RequireStatus(t, res, http.StatusNoContent)

		c, err := tu.app.queries.FindCommentById(com.Id)

		require.Empty(t, c)
		require.Error(t, err)
	})
}
