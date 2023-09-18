package app

import (
	"database/sql"
	"encoding/json"
	"io"
	"jira-clone/packages/queries"
	"jira-clone/packages/response"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		app.badRequest(w, "", nil)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.badRequest(w, "error reading req body", err)
		return
	}

	var projectDTO queries.CreateProjectDTO

	err = json.Unmarshal(body, &projectDTO)

	if err != nil {
		app.badRequest(w, "error decoding req body", err)
		return
	}

	err = projectDTO.Validate()

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	projectDTO.Created_by_id = r.Context().Value(ContextKey("userId")).(string)

	project, err := app.queries.CreateProject(projectDTO)

	if err != nil {
		app.badRequest(w, "", err)
		return
	}

	_, err = app.queries.CreateUserProjectXref(projectDTO.Created_by_id, project.Id)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusCreated, *project)
}

func (app *application) getProjectsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(ContextKey("userId")).(string)

	projects, err := app.queries.SelectProjectsForUser(userId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, projects)
}

func (app *application) getProjectDetails(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["projectId"]

	project, err := app.queries.GetJoinedProjectDetails(projectId)

	if err != nil {
		app.serverError(w, "", err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, project)
}

func (app *application) updateProject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.badRequest(w, "error reading req body", err)
		return
	}

	var payload queries.UpdateProjectDTO

	err = json.Unmarshal(body, &payload)

	if err != nil {
		app.badRequest(w, "error decoding req body", err)
		return
	}

	projectId := mux.Vars(r)["projectId"]

	project, err := app.queries.UpdateProject(projectId, payload)

	if err != nil {
		app.serverError(w, "", err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, project)
}

func (app *application) getProjectMembers(w http.ResponseWriter, r *http.Request) {
	projectId := mux.Vars(r)["projectId"]

	_, err := app.queries.GetProjectDetails(projectId)

	if err != nil {
		if err == sql.ErrNoRows {
			app.notFound(w, "project with id "+projectId+" not found", err)
			return
		}
		app.serverError(w, err.Error(), err)
		return
	}

	users, err := app.queries.GetProjectMembers(projectId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	response.NewSuccessResponse(w, http.StatusOK, users)
}
