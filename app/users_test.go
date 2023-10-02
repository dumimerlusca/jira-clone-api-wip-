package app

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserDetailsHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodGet, "/api/users/details/232dwdsa", struct{}{})

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})
	t.Run("should return 404 status if user does not exist", func(t *testing.T) {
		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/users/details/asjdnausdasudasduasd", struct{}{}, "")

		tu.RequireStatus(t, res, http.StatusNotFound)
	})

	t.Run("should return 200 status and the user details", func(t *testing.T) {
		user, _ := tu.app.queries.CreateRandomUser(t)

		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/users/details/"+user.Username, struct{}{}, user.Id)

		tu.RequireStatus(t, res, http.StatusOK)
		tu.IsSuccessResponseWithData(t, res)
	})
}

func TestGetWorkspaceMembers(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodGet, "/api/users/workspace-members", struct{}{})

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should return 200 status and a list of project members from all projects the current user is member off", func(t *testing.T) {
		p, u := tu.app.queries.CreateRandomProject(t)
		u1, _ := tu.app.queries.CreateRandomUser(t)
		tu.app.queries.CreateUserProjectXref(u1.Id, p.Id)
		u2, _ := tu.app.queries.CreateRandomUser(t)
		tu.app.queries.CreateUserProjectXref(u2.Id, p.Id)
		// P1 should have 3 members (u, u1, u2)

		// Just for adding more data
		tu.app.queries.CreateRandomProject(t)

		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/users/workspace-members", struct{}{}, u.Id)

		tu.RequireStatus(t, res, http.StatusOK)

		data := tu.GetSuccessResponseData(t, res).([]any)
		assert.Equal(t, 3, len(data))

		first := data[0].(map[string]any)

		assert.NotEmpty(t, first["id"])
		assert.NotEmpty(t, first["username"])
	})

	t.Run("should return project members if query parameter projectId is provided", func(t *testing.T) {
		p1, u := tu.app.queries.CreateRandomProject(t)
		p2 := tu.app.queries.CreateRandomProjectForUser(t, u.Id)

		u1, _ := tu.app.queries.CreateRandomUser(t)
		tu.app.queries.CreateUserProjectXref(u1.Id, p1.Id)
		tu.app.queries.CreateUserProjectXref(u1.Id, p2.Id)
		u2, _ := tu.app.queries.CreateRandomUser(t)
		tu.app.queries.CreateUserProjectXref(u2.Id, p1.Id)
		// P1 should have 3 members (u, u1, u2)
		// P2 should have 2 members (u, u1)

		// Just for adding more data
		tu.app.queries.CreateRandomProject(t)

		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/users/workspace-members?projectId="+p2.Id, struct{}{}, u.Id)

		tu.RequireStatus(t, res, http.StatusOK)

		data := tu.GetSuccessResponseData(t, res).([]any)
		assert.Equal(t, 2, len(data))
	})

	t.Run("should return 401 status if user is not project member", func(t *testing.T) {
		p, _ := tu.app.queries.CreateRandomProject(t)

		u1, _ := tu.app.queries.CreateRandomUser(t)

		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/users/workspace-members?projectId="+p.Id, struct{}{}, u1.Id)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})
}
