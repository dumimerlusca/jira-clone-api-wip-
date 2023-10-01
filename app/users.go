package app

import (
	"database/sql"
	"fmt"
	"jira-clone/packages/queries"
	"jira-clone/packages/response"
	"jira-clone/packages/util"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
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

func (app *application) getWorkspaceMembers(w http.ResponseWriter, r *http.Request) {
	userId := extractUserId(r)
	projectId := util.GetQueryParameter("projectId", r)

	projectIds, err := app.queries.GetProjectIdsWhereUserIsMember(userId)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	if len(projectIds) == 0 {
		response.NewSuccessResponse(w, http.StatusOK, []string{})
		return
	}

	if projectId != "" && !slices.Contains(projectIds, projectId) {
		app.unauthorizedRequest(w, "you are not member of this project", nil)
		return
	}

	mappedIds := []string{}
	for _, str := range projectIds {
		mappedIds = append(mappedIds, "'"+str+"'")
	}

	sqlWhere := ``

	if projectId != "" {
		sqlWhere = `WHERE project_id=` + fmt.Sprintf(`'%v'`, projectId)
	} else {
		sqlWhere = `WHERE project_id IN ` + fmt.Sprintf("(%v)", strings.Join(mappedIds, ","))
	}

	sql := `SELECT u.id, u.username FROM user_project_xref INNER JOIN users as u on user_id=u.id ` + sqlWhere + ` GROUP BY u.id`

	fmt.Println(sql)

	rows, err := app.db.Query(sql)

	if err != nil {
		app.serverError(w, err.Error(), err)
		return
	}

	users := []*queries.UserItem{}

	for rows.Next() {
		var user queries.UserItem

		err := rows.Scan(&user.Id, &user.Username)

		if err != nil {
			app.serverError(w, err.Error(), err)
			return
		}

		users = append(users, &user)
	}

	response.NewSuccessResponse(w, http.StatusOK, users)
}
