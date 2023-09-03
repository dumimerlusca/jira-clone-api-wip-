package app

import (
	"context"
	"jira-clone/packages/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestAuthMW(t *testing.T) {
	t.Run("should return 401 if token is not valid and should not call the handler", func(t *testing.T) {
		mockHandler := mock.NewMockHandler()
		handlerFunc := tApp.authMW(mockHandler.HandleFunc)

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "", nil)

		req.Header.Set("Authorization", "Bearer ajsdnasd")

		handlerFunc(res, req)

		require.Equal(t, http.StatusUnauthorized, res.Code)
		require.Equal(t, 0, mockHandler.CallsCount)
	})

	t.Run("should call the handler if token is valid and add userId to req context", func(t *testing.T) {

		user, _ := tApp.queries.CreateRandomUser(t)
		token, _ := generateAuthToken(user.Id, user.Username)

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "", nil)

		req.Header.Set("Authorization", "Bearer"+" "+token)

		tApp.authMW(authMockHandler(t, user.Id))(res, req)

		require.Equal(t, http.StatusOK, res.Code)
	})
}

func authMockHandler(t *testing.T, expectedUserId string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(ContextKey("userId"))

		require.Equal(t, expectedUserId, userId)

		w.WriteHeader(http.StatusOK)
	}
}

func TestIsProjectOwnerMW(t *testing.T) {
	t.Run("should verify that the current logged in user is project owner", func(t *testing.T) {
		mockHandler := mock.NewMockHandler()
		project, user := tApp.queries.CreateRandomProject(t)

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "", nil)

		ctx := context.WithValue(req.Context(), ContextKey("userId"), user.Id)

		req = req.WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"projectId": project.Id})

		tApp.isProjectOwnerMW(mockHandler.HandleFunc)(res, req)

		require.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("should return 401 if current logged in user is not project owner", func(t *testing.T) {
		mockHandler := mock.NewMockHandler()
		project, _ := tApp.queries.CreateRandomProject(t)

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "", nil)

		ctx := context.WithValue(req.Context(), ContextKey("userId"), "otherUserId")

		req = req.WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"projectId": project.Id})

		tApp.isProjectOwnerMW(mockHandler.HandleFunc)(res, req)

		require.Equal(t, http.StatusUnauthorized, res.Code)
	})
}

func TestIsProjectMemberMW(t *testing.T) {
	t.Run("should work with ticketId as url parameter", func(t *testing.T) {
		mockHandler := mock.NewMockHandler()

		res, req := tu.CreateReqAndRes(t, "POST", "", nil)

		ticket := tu.app.queries.CreateRandomTicket(t)

		ctx := context.WithValue(req.Context(), ContextKey("userId"), ticket.Created_by_id)

		req = req.WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"ticketId": ticket.Id})

		tu.app.isProjectMemberMW(mockHandler.HandleFunc)(res, req)

		require.Equal(t, 1, mockHandler.CallsCount)
	})

	t.Run("should work with projectId as url parameter", func(t *testing.T) {
		mockHandler := mock.NewMockHandler()

		res, req := tu.CreateReqAndRes(t, "POST", "", nil)

		project, _ := tu.app.queries.CreateRandomProject(t)

		ctx := context.WithValue(req.Context(), ContextKey("userId"), project.Created_by_id)

		req = req.WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"projectId": project.Id})

		tu.app.isProjectMemberMW(mockHandler.HandleFunc)(res, req)

		require.Equal(t, 1, mockHandler.CallsCount)
	})

	t.Run("should return 401 if user is not project member", func(t *testing.T) {
		mockHandler := mock.NewMockHandler()

		res, req := tu.CreateReqAndRes(t, "POST", "", nil)

		project, _ := tu.app.queries.CreateRandomProject(t)

		user, _ := tu.app.queries.CreateRandomUser(t)

		ctx := context.WithValue(req.Context(), ContextKey("userId"), user.Id)

		req = req.WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"projectId": project.Id})

		tu.app.isProjectMemberMW(mockHandler.HandleFunc)(res, req)

		require.Equal(t, 0, mockHandler.CallsCount)
		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})
	t.Run("should return 400 bad request if projectId or ticket id are missing", func(t *testing.T) {
		mockHandler := mock.NewMockHandler()
		res, req := tu.CreateReqAndRes(t, "POST", "", nil)
		ctx := context.WithValue(req.Context(), ContextKey("userId"), "someone")
		req = req.WithContext(ctx)
		tu.app.isProjectMemberMW(mockHandler.HandleFunc)(res, req)
		require.Equal(t, 0, mockHandler.CallsCount)
		tu.RequireStatus(t, res, http.StatusBadRequest)
	})
}
