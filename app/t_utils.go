package app

import (
	"bytes"
	"encoding/json"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestUtils struct {
	app *application
}

func (u *TestUtils) IsSuccessResponseWithData(t *testing.T, res *httptest.ResponseRecorder) {
	var a response.SuccessResponse

	err := util.ReadAndUnmarshal(res.Body, &a)

	require.NoError(t, err)

	require.Equal(t, true, a.Success)
	require.NotZero(t, a.Data)
}

func (u *TestUtils) RequireStatus(t *testing.T, res *httptest.ResponseRecorder, status int) {
	require.Equal(t, status, res.Code)
}

func (u *TestUtils) AuthorizeRequest(t *testing.T, req *http.Request) {
	user, _ := u.app.queries.CreateRandomUser(t)
	token, err := generateAuthToken(user.Id, user.Username)
	require.NoError(t, err)
	req.Header.Set("Authorization", `Bearer`+" "+token)
}

func (u *TestUtils) CreateReqAndRes(t *testing.T, method string, url string, body any) (*httptest.ResponseRecorder, *http.Request) {
	d, err := json.Marshal(body)

	require.NoError(t, err)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(d))

	require.NoError(t, err)

	res := httptest.NewRecorder()

	return res, req
}

func (u *TestUtils) CreateAuthorizedReqAndRes(t *testing.T, method string, url string, body any) (*httptest.ResponseRecorder, *http.Request) {
	res, req := u.CreateReqAndRes(t, method, url, body)

	u.AuthorizeRequest(t, req)

	return res, req
}

func (u *TestUtils) SendAuthorizedReq(t *testing.T, method string, url string, body any) (*httptest.ResponseRecorder, *http.Request) {
	res, req := u.CreateAuthorizedReqAndRes(t, method, url, body)
	u.app.routes().ServeHTTP(res, req)
	return res, req
}

func (u *TestUtils) SendUnauthorizedReq(t *testing.T, method string, url string, body any) (*httptest.ResponseRecorder, *http.Request) {
	res, req := u.CreateReqAndRes(t, method, url, body)
	u.app.routes().ServeHTTP(res, req)
	return res, req
}

func (u *TestUtils) RequireAuth(t *testing.T, res *httptest.ResponseRecorder) {

}
