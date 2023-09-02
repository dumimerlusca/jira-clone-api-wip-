package app

import (
	"database/sql"
	"jira-clone/packages/response"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) getUserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["username"]

	user, err := app.queries.FindUserByUsername(userId, false)

	if err != nil {
		if err == sql.ErrNoRows {
			app.notFound(w, "user with id "+userId+" not found", err)
			return
		}
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, user)
}
