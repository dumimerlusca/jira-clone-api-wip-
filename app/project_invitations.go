package app

import (
	"context"
	"jira-clone/packages/consts"
	"jira-clone/packages/queries"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"

	"github.com/gorilla/mux"
)

type sendProjectInvitationReqPayload struct {
	ReceiverId string
}

func (app *application) sendProjectInvitationHandler(w http.ResponseWriter, r *http.Request) {
	sentBy := r.Context().Value(ContextKey("userId"))
	projectId := mux.Vars(r)["projectId"]

	var payload sendProjectInvitationReqPayload

	err := util.ReadAndUnmarshal(r.Body, &payload)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	if payload.ReceiverId == "" {
		app.badRequest(w, "receiver id must not be empty", nil)
		return
	}

	count, err := app.queries.SelectPendingProjectInvitationsCount(payload.ReceiverId, projectId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	// There should not exist more than one pending invitation
	if *count != 0 {
		app.badRequest(w, consts.MsgOneProjectInvitationIsAlreadyPending, err)
		return

	}

	_, err = app.queries.CreateProjectInvitation(queries.CreateProjectInvitationPayload{Receiver_id: payload.ReceiverId, Project_id: projectId, Sender_id: sentBy.(string), Status: "pending"})

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusCreated, nil)
}

func (app *application) acceptProjectInviteHandler(w http.ResponseWriter, r *http.Request) {
	inviteId := mux.Vars(r)["inviteId"]
	userId := r.Context().Value(ContextKey("userId"))
	projectId := r.Context().Value(ContextKey("projectId"))

	err := app.queries.UpdateProjectInvitationStatus(inviteId, "accepted")

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	_, err = app.queries.CreateUserProjectXref(userId.(string), projectId.(string))

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, nil)
}

func (app *application) rejectProjectInviteHandler(w http.ResponseWriter, r *http.Request) {
	inviteId := mux.Vars(r)["inviteId"]
	err := app.queries.UpdateProjectInvitationStatus(inviteId, "rejected")

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, nil)
}

func (app *application) validateUpdateProjectInviteMW(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(ContextKey("userId"))
		inviteId := mux.Vars(r)["inviteId"]

		if inviteId == "" || userId == "" {
			app.badRequest(w, "", nil)
			return
		}

		inv, err := app.queries.FindProjectInvitationById(inviteId)

		if err != nil {
			app.serverError(w, err.Error(), err)
			return
		}

		if userId != inv.Receiver_id {
			app.unauthorizedRequest(w, "", nil)
			return
		}

		if inv.Status != "pending" {
			app.badRequest(w, consts.MsgProjectInvitationNotInPendingState, nil)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey("projectId"), inv.Project_id)

		r = r.WithContext(ctx)

		handler(w, r)
	}

}

func (app *application) getSentProjectInvites(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(ContextKey("userId")).(string)

	invitations, err := app.queries.SelectSentProjectInvites(userId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, invitations)

}

func (app *application) getReceivedProjectInvites(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(ContextKey("userId")).(string)

	invitations, err := app.queries.SelectReceivedProjectInvites(userId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, invitations)

}
