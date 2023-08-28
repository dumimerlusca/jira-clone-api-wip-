package app

import (
	"fmt"
	"jira-clone/packages/response"
	"net/http"
)

func (app *application) errorMessage(w http.ResponseWriter, status int, message string, err error) {
	fmt.Println(message)
	if err != nil {
		fmt.Println(err.Error())
	}

	response.NewErrorResponse(w, status, message)
}

func (app *application) serverError(w http.ResponseWriter, message string, err error) {
	defMsg := "The server encountered a problem and could not process your request"

	if message == "" {
		message = defMsg
	}

	app.errorMessage(w, http.StatusInternalServerError, message, err)
}

func (app *application) badRequest(w http.ResponseWriter, message string, err error) {
	defMsg := "Bad request"

	if message == "" {
		message = defMsg
	}

	app.errorMessage(w, http.StatusBadRequest, message, err)
}

func (app *application) unauthorizedRequest(w http.ResponseWriter, message string, err error) {
	defMsg := "Unauthorized"
	if message == "" {
		message = defMsg
	}
	app.errorMessage(w, http.StatusUnauthorized, message, err)
}
