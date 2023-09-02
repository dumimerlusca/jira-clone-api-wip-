package app

import (
	"context"
	"jira-clone/packages/consts"
	"jira-clone/packages/mock"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendProjectInvitation(t *testing.T) {
	t.Run("should properly create and save invitation to the database", func(t *testing.T) {
		project, sender := tApp.queries.CreateRandomProject(t)
		receiver, _ := tApp.queries.CreateRandomUser(t)

		payload := sendProjectInvitationReqPayload{ReceiverId: receiver.Id}

		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, "/api/projects/sendInvite/"+project.Id, payload, sender.Id)

		tu.RequireStatus(t, res, http.StatusCreated)

		var count int

		row := tApp.db.QueryRow(`SELECT COUNT(*) FROM project_invitations WHERE receiver_id=$1 AND project_id=$2 AND sender_id=$3 AND status='pending'`, receiver.Id, project.Id, sender.Id)

		err := row.Scan(&count)

		require.NoError(t, err)
		require.Equal(t, 1, count)

	})

	t.Run("should not create new invitation if one is already pending", func(t *testing.T) {
		existentInvite, err := tApp.queries.CreateRandomProjectInvite(t, "pending")

		require.NoError(t, err)

		payload := sendProjectInvitationReqPayload{ReceiverId: existentInvite.Receiver.Id}

		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, "/api/projects/sendInvite/"+existentInvite.Project.Id, payload, existentInvite.Sender.Id)

		var resEror response.ErrorResponse

		util.ReadAndUnmarshal(res.Body, &resEror)

		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, consts.MsgOneProjectInvitationIsAlreadyPending, resEror.Error)

	})
}

func newUpdateProjectInviteRequest(loggedUserId string, inviteId string, handleFunc http.HandlerFunc) (*httptest.ResponseRecorder, *http.Request) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "", nil)

	ctx := context.WithValue(req.Context(), ContextKey("userId"), loggedUserId)

	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"inviteId": inviteId})

	handler := tApp.validateUpdateProjectInviteMW(handleFunc)

	handler(res, req)

	return res, req
}

func TestAcceptProjectInviteHandler(t *testing.T) {
	t.Run("should update invite status to 'accepted'", func(t *testing.T) {
		inv, err := tApp.queries.CreateRandomProjectInvite(t, "pending")

		require.NoError(t, err)

		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, "/api/projects/acceptInvite/"+inv.Inv.Id, nil, inv.Inv.Receiver_id)

		require.Equal(t, http.StatusOK, res.Code)

		var successRes response.SuccessResponse
		util.ReadAndUnmarshal(res.Body, &successRes)
		require.Equal(t, true, successRes.Success)

		updatedIv, err := tApp.queries.FindProjectInvitationById(inv.Inv.Id)

		assert.NoError(t, err)
		assert.Equal(t, "accepted", updatedIv.Status)

		// Should create user_project_xref
		count, err := tApp.queries.IsUserInProject(inv.Inv.Receiver_id, inv.Inv.Project_id)

		assert.NoError(t, err)
		assert.Equal(t, 1, count)

	})

}
func TestRejectProjectInviteHandler(t *testing.T) {
	t.Run("should update invite status to 'rejected'", func(t *testing.T) {
		inv, err := tApp.queries.CreateRandomProjectInvite(t, "pending")

		require.NoError(t, err)

		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, "/api/projects/rejectInvite/"+inv.Inv.Id, nil, inv.Inv.Receiver_id)

		require.Equal(t, http.StatusOK, res.Code)

		var successRes response.SuccessResponse
		util.ReadAndUnmarshal(res.Body, &successRes)
		require.Equal(t, true, successRes.Success)

		updatedIv, err := tApp.queries.FindProjectInvitationById(inv.Inv.Id)

		require.NoError(t, err)

		require.Equal(t, "rejected", updatedIv.Status)

		// Should not create user_project_xref
		count, err := tApp.queries.IsUserInProject(inv.Inv.Receiver_id, inv.Inv.Project_id)

		assert.NoError(t, err)
		assert.Equal(t, 0, count)

	})

}

func TestValidateUpdateProjectInviteMW(t *testing.T) {
	t.Run("should return 401 if another user is trying to update the status", func(t *testing.T) {
		inv, _ := tApp.queries.CreateRandomProjectInvite(t, "pending")
		randomUser, _ := tApp.queries.CreateRandomUser(t)

		mockHandler := mock.NewMockHandler()

		res, _ := newUpdateProjectInviteRequest(randomUser.Id, inv.Inv.Id, mockHandler.HandleFunc)

		assert.Equal(t, 0, mockHandler.CallsCount)
		require.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("should return 400 if the invitation is in any other state other than 'pending'", func(t *testing.T) {
		inv, _ := tApp.queries.CreateRandomProjectInvite(t, "rejected")

		mockHandler := mock.NewMockHandler()

		res, _ := newUpdateProjectInviteRequest(inv.Receiver.Id, inv.Inv.Id, mockHandler.HandleFunc)

		var resErr response.ErrorResponse
		util.ReadAndUnmarshal(res.Body, &resErr)

		assert.Equal(t, 0, mockHandler.CallsCount)
		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, resErr.Error, consts.MsgProjectInvitationNotInPendingState)
	})
}
