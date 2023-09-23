package app

import (
	"jira-clone/packages/queries"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"

	"github.com/gorilla/mux"
)

type createCommentRequestPayload struct {
	Text string
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(ContextKey("userId")).(string)
	ticketId := mux.Vars(r)["ticketId"]

	var payload createCommentRequestPayload

	err := util.ReadAndUnmarshal(r.Body, &payload)

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	if payload.Text == "" {
		app.badRequest(w, "text must not be empty", nil)
		return
	}

	coment, err := app.queries.CreateComment(queries.CreateCommentDTO{Author_id: userId, Ticket_id: ticketId, Text: payload.Text})

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusCreated, coment)
}

func (app *application) getTicketCommentsHandler(w http.ResponseWriter, r *http.Request) {
	ticketId := mux.Vars(r)["ticketId"]

	tickets, err := app.queries.SelectJoinedTicketComments(ticketId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, tickets)
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentId := mux.Vars(r)["commentId"]
	userId := r.Context().Value(ContextKey("userId"))

	comment, err := app.queries.FindCommentById(commentId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	if userId != comment.Author_id {
		app.unauthorizedRequest(w, "you can't delete this comment", err)
		return
	}

	err = app.queries.DeleteComment(commentId)

	if err != nil {
		app.serverError(w, "comment deletion failed", err)
		return
	}

	response.NewSuccessResponse(w, http.StatusNoContent, nil)
}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentId := mux.Vars(r)["commentId"]
	userId := r.Context().Value(ContextKey("userId"))

	comment, err := app.queries.FindCommentById(commentId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	if userId != comment.Author_id {
		app.unauthorizedRequest(w, "you can't delete this comment", err)
		return
	}

	var payload queries.UpdateCommentPayload

	err = util.ReadAndUnmarshal(r.Body, &payload)

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	err = payload.Validate()

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	updatedCom, err := app.queries.UpdateComment(commentId, payload)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, updatedCom)
}
