package app

import (
	"jira-clone/packages/consts"
	"jira-clone/packages/queries"
	"jira-clone/packages/random"
	"net/http"
	"testing"
)

func TestCreateProjectHandler(t *testing.T) {
	t.Run("should require auth", func(t *testing.T) {
		res, _ := tu.SendUnauthorizedReq(t, http.MethodPost, consts.ApiPathCreateProject, struct{}{})
		tu.RequireStatus(t, res, http.StatusUnauthorized)
	})

	t.Run("create project and save it to database", func(t *testing.T) {
		payload := queries.CreateProjectDTO{Name: random.RandomString(10), Key: random.RandomString(4), Description: random.RandomString(30)}

		res, _ := tu.SendAuthorizedReq(t, http.MethodPost, consts.ApiPathCreateProject, payload)

		tu.RequireStatus(t, res, http.StatusCreated)
		tu.IsSuccessResponseWithData(t, res)
	})
}
