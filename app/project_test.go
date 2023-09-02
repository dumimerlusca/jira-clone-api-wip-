package app

import (
	"jira-clone/packages/queries"
	"jira-clone/packages/random"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateProjectHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodPost, "/api/projects/create", struct{}{})
		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("create project and save it to database", func(t *testing.T) {
		payload := queries.CreateProjectDTO{Name: random.RandomString(10), Key: random.RandomString(4), Description: random.RandomString(30)}

		user, _ := tu.app.queries.CreateRandomUser(t)

		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, "/api/projects/create", payload, user.Id)

		tu.RequireStatus(t, res, http.StatusCreated)
		var resBody response.SuccessResponse

		util.ReadAndUnmarshal(res.Body, &resBody)

		data, ok := resBody.Data.(map[string]any)

		require.Equal(t, true, ok)

		projectId := data["id"].(string)

		require.NotEmpty(t, projectId)
		isProjectMember, err := tu.app.queries.IsProjectMember(user.Id, projectId)

		require.NoError(t, err)
		require.Equal(t, true, isProjectMember)

	})
}

func TestGetProjectDetailsHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodGet, "/api/projects/details/2732462424", struct{}{})
		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should return project details", func(t *testing.T) {
		project, user := tApp.queries.CreateRandomProject(t)

		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/projects/details/"+project.Id, nil, user.Id)

		tu.RequireStatus(t, res, http.StatusOK)
		tu.IsSuccessResponseWithData(t, res)

	})
}

func TestUpdateProjectHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodPatch, "/api/projects/update/2732462424", struct{}{})
		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("should update project", func(t *testing.T) {
		project, user := tApp.queries.CreateRandomProject(t)

		name := "NEW NAME"
		desc := "NEW DESCRIPTION"
		key := "KEY1"

		payload := queries.UpdateProjectDTO{Name: &name, Description: &desc, Key: &key}

		res, _ := tu.SendAuthorizedReq(t, http.MethodPatch, "/api/projects/update/"+project.Id, payload, user.Id)

		tu.RequireStatus(t, res, http.StatusOK)
		tu.IsSuccessResponseWithData(t, res)

		updatedProject, err := tApp.queries.GetProjectDetails(project.Id)

		require.NoError(t, err)

		require.Equal(t, name, updatedProject.Name)
		require.Equal(t, desc, updatedProject.Description)
		require.Equal(t, key, updatedProject.Key)

	})
}

func TestGetProjectMembers(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodGet, "/api/projects/users/someprojectId", struct{}{})

		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})
	t.Run("should return 404 if project does not exist", func(t *testing.T) {
		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/projects/users/someprojectIdasdasd", struct{}{}, "")

		tu.RequireStatus(t, res, http.StatusNotFound)
	})

	t.Run("should return 200 and a list with all project members", func(t *testing.T) {
		project1, _ := tu.app.queries.CreateRandomProject(t)

		user1, _ := tu.app.queries.CreateRandomUser(t)
		user2, _ := tu.app.queries.CreateRandomUser(t)

		_, err := tu.app.db.Exec(`INSERT INTO user_project_xref(user_id, project_id) VALUES
		($1,$2), ($3,$4)`, user1.Id, project1.Id, user2.Id, project1.Id)

		require.NoError(t, err)

		res, _ := tu.SendAuthorizedReq(t, http.MethodGet, "/api/projects/users/"+project1.Id, nil, "")

		tu.RequireStatus(t, res, http.StatusOK)

		var resBody response.SuccessResponse

		util.ReadAndUnmarshal(res.Body, &resBody)

		resData, ok := resBody.Data.([]interface{})
		require.Equal(t, true, ok)

		// owner + the other 2 created above
		require.Equal(t, 3, len(resData))
	})
}
