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

func (u *TestUtils) IsErrorResponse(t *testing.T, res *httptest.ResponseRecorder) {
	var a response.ErrorResponse

	err := util.ReadAndUnmarshal(res.Body, &a)

	require.NoError(t, err)

	require.Equal(t, false, a.Success)
	require.NotZero(t, a.Error)
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

func (u *TestUtils) AuthorizeRequest(t *testing.T, req *http.Request, userId string) {
	if userId == "" {
		user, _ := u.app.queries.CreateRandomUser(t)
		userId = user.Id
	}

	token, err := generateAuthToken(userId, "")
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

func (u *TestUtils) CreateAuthorizedReqAndRes(t *testing.T, method string, url string, body any, userId string) (*httptest.ResponseRecorder, *http.Request) {
	res, req := u.CreateReqAndRes(t, method, url, body)

	u.AuthorizeRequest(t, req, userId)

	return res, req
}

func (u *TestUtils) SendAuthorizedReq(t *testing.T, method string, url string, body any, userId string) (*httptest.ResponseRecorder, *http.Request) {
	res, req := u.CreateAuthorizedReqAndRes(t, method, url, body, userId)
	u.app.routes().ServeHTTP(res, req)
	return res, req
}

func (u *TestUtils) SendUnauthorizedReq(t *testing.T, method string, url string, body any) (*httptest.ResponseRecorder, *http.Request) {
	res, req := u.CreateReqAndRes(t, method, url, body)
	u.app.routes().ServeHTTP(res, req)
	return res, req
}

func (U *TestUtils) GetSuccessResponseData(t *testing.T, res *httptest.ResponseRecorder) any {
	var body response.SuccessResponse

	util.ReadAndUnmarshal(res.Body, &body)

	require.Equal(t, true, body.Success)
	require.NotNil(t, body.Data)

	return body.Data
}
