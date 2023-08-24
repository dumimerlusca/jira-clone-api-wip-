package app

import (
	"encoding/json"
	"io"
	"jira-clone/packages/db"
	"jira-clone/packages/response"
	"net/http"

	"github.com/gorilla/mux"
)

// require auth
func (app *application) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	var projectDTO db.CreateProjectDTO

	err = json.Unmarshal(body, &projectDTO)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	err = projectDTO.Validate()

	if err != nil {
		app.badRequest(w, err)
		return
	}

	projectDTO.Leader_id = "1" // TODO Get user id from logged in user

	project, err := app.queries.CreateProject(projectDTO)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	response.JSONWithHeaders(w, http.StatusCreated, &project)
}

// require auth
func (app *application) getProjectsHandler(w http.ResponseWriter, r *http.Request) {

}

// require auth
func (app *application) getProjectDetails(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["projectId"]

	project, err := app.queries.GetProjectDetails(projectId)

	if err != nil {
		app.serverError(w, err)
		return
	}

	response.JSONWithHeaders(w, http.StatusOK, project)
}

// require auth
func (app *application) updateProject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	var payload db.UpdateProjectDTO

	err = json.Unmarshal(body, &payload)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	projectId := mux.Vars(r)["projectId"]

	project, err := app.queries.UpdateProject(projectId, payload)

	if err != nil {
		app.serverError(w, err)
		return
	}

	response.JSONWithHeaders(w, http.StatusOK, project)
}
