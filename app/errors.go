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

	err = response.JSONWithHeaders(w, status, map[string]string{"error": message})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) serverError(w http.ResponseWriter, err error) {
	message := "The server encountered a problem and could not process your request"

	app.errorMessage(w, http.StatusInternalServerError, message, err)
}

func (app *application) badRequest(w http.ResponseWriter, err error) {
	message := "Bad request"
	if err != nil {
		message = err.Error()
	}
	app.errorMessage(w, http.StatusBadRequest, message, err)
}

func (app *application) unauthorizedRequest(w http.ResponseWriter, err error) {
	msg := "Unauthorized"
	app.errorMessage(w, http.StatusUnauthorized, msg, err)
}
