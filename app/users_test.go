package app

import (
	"net/http"
	"testing"
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
