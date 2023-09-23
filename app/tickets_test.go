package app

import (
	"fmt"
	"jira-clone/packages/random"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTicketHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodPost, "/api/projects/1/tickets/create", nil)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("require user to be project member", func(t *testing.T) {
		project, _ := tu.app.queries.CreateRandomProject(t)
		user, _ := tu.app.queries.CreateRandomUser(t)
		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, "/api/projects/"+project.Id+"/tickets/create", nil, user.Id)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should save ticket to database and return it to client", func(t *testing.T) {
		project, user := tu.app.queries.CreateRandomProject(t)

		payload := struct {
			Assignee_id  string
			Story_points int
			Description  string
			Priority     int
			Title        string
			Type         string
		}{Assignee_id: user.Id, Story_points: 100, Description: random.RandomString(30), Title: random.RandomString(10), Priority: 0, Type: "epic"}

		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, "/api/projects/"+project.Id+"/tickets/create", payload, user.Id)

		tu.RequireStatus(t, res, http.StatusCreated)

		var resBody response.SuccessResponse

		util.ReadAndUnmarshal(res.Body, &resBody)

		data, ok := resBody.Data.(map[string]any)

		assert.Equal(t, true, ok)
		assert.Equal(t, project.Id, data["project_id"])
		assert.Equal(t, user.Id, data["created_by_id"])
		assert.Equal(t, payload.Title, data["title"])
		assert.Equal(t, payload.Type, data["type"])
		assert.Equal(t, payload.Description, data["description"])
		assert.Equal(t, payload.Assignee_id, data["assignee_id"])
		assert.Equal(t, float64(payload.Priority), data["priority"])
		assert.Equal(t, "open", data["status"])

	})
}

func TestUpdateTicketHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodPatch, "/api/tickets/update/ticketId", nil)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("require user to be project member", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)
		user, _ := tu.app.queries.CreateRandomUser(t)
		res, _ := tu.SendAuthorizedReq(t, http.MethodPatch, "/api/tickets/update/"+ticket.Id, nil, user.Id)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should save updated fields to database and return it to client", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)

		user, _ := tu.app.queries.CreateRandomUser(t)

		payload := struct {
			Assignee_id  string
			Story_points int
			Description  string
			Priority     int
			Title        string
			Status       string
		}{Assignee_id: user.Id, Story_points: 100, Description: random.RandomString(30), Title: random.RandomString(10), Priority: 0, Status: "tested"}

		res, _ := tu.SendAuthorizedReq(t, http.MethodPatch, "/api/tickets/update/"+ticket.Id, payload, ticket.Created_by_id)

		tu.RequireStatus(t, res, http.StatusOK)

		var resBody response.SuccessResponse

		util.ReadAndUnmarshal(res.Body, &resBody)

		data, ok := resBody.Data.(map[string]any)

		assert.Equal(t, true, ok)
		assert.Equal(t, payload.Title, data["title"])
		assert.Equal(t, payload.Description, data["description"])
		assert.Equal(t, payload.Assignee_id, data["assignee_id"])
		assert.Equal(t, float64(payload.Story_points), data["story_points"])
		assert.Equal(t, float64(payload.Priority), data["priority"])
		assert.Equal(t, "tested", data["status"])

	})
}

func TestGetProjectTickets(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodGet, "/api/projects/projectId/tickets", nil)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("require user to be project member", func(t *testing.T) {
		project, _ := tu.app.queries.CreateRandomProject(t)
		user, _ := tu.app.queries.CreateRandomUser(t)
		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/projects/"+project.Id+"/tickets", nil, user.Id)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should return 200 and the array of tickets", func(t *testing.T) {

		ticket := tu.app.queries.CreateRandomTicket(t)

		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/projects/"+ticket.Project_id+"/tickets", nil, ticket.Created_by_id)

		tu.RequireStatus(t, res, http.StatusOK)

		var resBody response.SuccessResponse

		util.ReadAndUnmarshal(res.Body, &resBody)

		data, ok := resBody.Data.([]any)

		assert.Equal(t, true, ok)
		assert.Equal(t, 1, len(data))
	})

}

func TestGetTicketDetailsHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodGet, "/api/tickets/VISA-1", nil)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should return 200 status and ticket details", func(t *testing.T) {
		ticket := tu.app.queries.CreateRandomTicket(t)
		key, err := tu.app.queries.FindTicketKeyById(ticket.Id)

		require.NoError(t, err)

		res, _ := tu.SendAuthorizedReq(t, "GET", fmt.Sprintf("/api/tickets/%v", *key), nil, ticket.Created_by_id)

		assert.Equal(t, http.StatusOK, res.Code)

		var resBody response.SuccessResponse

		util.ReadAndUnmarshal(res.Body, &resBody)

		data, ok := resBody.Data.(map[string]any)

		assert.Equal(t, true, ok)
		assert.Equal(t, ticket.Id, data["id"])
		assert.Equal(t, *key, data["key"])
	})
}
