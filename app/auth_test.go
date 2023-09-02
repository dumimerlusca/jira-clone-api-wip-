package app

import (
	"bytes"
	"encoding/json"
	"io"
	"jira-clone/packages/random"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newRegisterRequest(t *testing.T, body any) (*httptest.ResponseRecorder, *http.Request) {
	u, err := json.Marshal(body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(u))
	tApp.registerHandler(res, req)
	return res, req
}

func TestRegister(t *testing.T) {
	t.Run("should register the user successfully and return login token", func(t *testing.T) {
		p := registerRequestPayload{Username: random.RandomString(20), Password: random.RandomString(10)}

		res, _ := newRegisterRequest(t, p)

		require.Equal(t, http.StatusCreated, res.Code)

		user, err := tApp.queries.FindUserByUsername(p.Username, true)

		require.NoError(t, err)
		// Password should be hashed
		require.NotEqual(t, p.Password, user.Password)

		// Assert response body
		rData, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var rPayload registerResponsePayload

		json.Unmarshal(rData, &rPayload)

		require.NotZero(t, rPayload.Token)
		require.NotZero(t, rPayload.User)
		require.NotZero(t, rPayload.User.Id)
		require.NotZero(t, rPayload.User.Username)
		require.Zero(t, rPayload.User.Password)
	})

	t.Run("should return 400 err if bad values are provided", func(t *testing.T) {

		// IF wrong json is provided
		res, _ := newRegisterRequest(t, "bad value")
		require.Equal(t, http.StatusBadRequest, res.Code)

		// IF nil body
		res, _ = newRegisterRequest(t, nil)
		require.Equal(t, http.StatusBadRequest, res.Code)

		// IF password is missing
		payload := registerRequestPayload{Username: random.RandomString(10), Password: ""}
		res, _ = newRegisterRequest(t, payload)
		require.Equal(t, http.StatusBadRequest, res.Code)

		// IF username is missing
		payload = registerRequestPayload{Password: random.RandomString(10), Username: ""}
		res, _ = newRegisterRequest(t, payload)
		require.Equal(t, http.StatusBadRequest, res.Code)

		// IF password is less than 6 characters long
		payload = registerRequestPayload{Username: random.RandomString(10), Password: "12345"}
		res, _ = newRegisterRequest(t, payload)
		require.Equal(t, http.StatusBadRequest, res.Code)
	})

}

func newLoginRequest(t *testing.T, body any) (*httptest.ResponseRecorder, *http.Request) {
	u, err := json.Marshal(body)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(u))
	tApp.loginHandler(res, req)
	return res, req
}

func TestLogin(t *testing.T) {
	t.Run("should login the user successfuly", func(t *testing.T) {
		user := registerRequestPayload{Username: random.RandomString(20), Password: random.RandomString(10)}
		newRegisterRequest(t, user)

		payload := loginRequestPayload{Username: user.Username, Password: user.Password}
		res, _ := newLoginRequest(t, payload)

		require.Equal(t, http.StatusOK, res.Code)

		var rBody loginResponsePayload

		util.ReadAndUnmarshal(res.Body, &rBody)

		require.NotZero(t, rBody.Token)

	})
}

func TestGetCurrentUserHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, "GET", "/api/auth/me", nil)

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should return status 200 and current logged in user", func(t *testing.T) {
		user, _ := tu.app.queries.CreateRandomUser(t)

		res, _ := tu.SendAuthorizedReq(t, "GET", "/api/auth/me", nil, user.Id)

		tu.RequireStatus(t, res, http.StatusOK)

		var resBody response.SuccessResponse

		util.ReadAndUnmarshal(res.Body, &resBody)

		data, ok := resBody.Data.(map[string]any)

		require.Equal(t, true, ok)
		require.Equal(t, user.Id, data["id"])
		require.Equal(t, user.Username, data["username"])
		require.Equal(t, user.Created_at, data["created_at"])

	})
}
