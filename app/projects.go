package app

import (
	"encoding/json"
	"io"
	"jira-clone/packages/queries"
	"jira-clone/packages/response"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	var projectDTO queries.CreateProjectDTO

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

	projectDTO.Created_by_id = r.Context().Value(ContextKey("userId")).(string)

	project, err := app.queries.CreateProject(projectDTO)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	response.JSONWithHeaders(w, http.StatusCreated, &project)
}

func (app *application) getProjectsHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getProjectDetails(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["projectId"]

	project, err := app.queries.GetJoinedProjectDetails(projectId)

	if err != nil {
		app.serverError(w, err)
		return
	}

	response.JSONWithHeaders(w, http.StatusOK, project)
}

func (app *application) updateProject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	var payload queries.UpdateProjectDTO

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
