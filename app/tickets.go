package app

import (
	"net/http"
)

func (app *application) createTicketHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Create bug"))
}

func (app *application) updateTicketHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update bug"))
}
