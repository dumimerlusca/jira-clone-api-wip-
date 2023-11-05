package queries

import (
	"fmt"
	"jira-clone/packages/models"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type CreateProjectDTO struct {
	Name          string `json:"name"`
	Key           string `json:"key"`
	Description   string `json:"description"`
	Created_by_id string `json:"created_by_id"`
}

func (p *CreateProjectDTO) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name required")
	}

	if p.Key == "" {
		return fmt.Errorf("key required")
	}

	return nil
}

func (q *Queries) CreateProject(data CreateProjectDTO) (*models.Project, error) {
	id := uuid.NewString()

	sql := `INSERT INTO projects(id, name, key, description, created_by_id) VALUES($1, $2, $3, $4, $5) RETURNING id, name, key, description, created_by_id, created_at`

	row := q.Db.QueryRow(sql, id, data.Name, data.Key, data.Description, data.Created_by_id)

	var project models.Project

	err := row.Scan(&project.Id, &project.Name, &project.Key, &project.Description, &project.Created_by_id, &project.Created_at)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

type UpdateProjectDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Key         *string `json:"key"`
}

func (q *Queries) UpdateProject(projectId string, data UpdateProjectDTO) (*models.Project, error) {
	if data.Name == nil && data.Key == nil && data.Description == nil {
		return nil, fmt.Errorf("no fields provided")
	}

	values := []any{}
	var sqlColumnValues []string

	handleField := func(name string, value string) {
		values = append(values, value)
		s := name + "=$" + strconv.FormatInt(int64(len(values)), 10)
		sqlColumnValues = append(sqlColumnValues, s)
	}

	if data.Name != nil {
		handleField("name", *data.Name)
	}

	if data.Description != nil {
		handleField("description", *data.Description)
	}

	if data.Key != nil {
		handleField("key", *data.Key)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("no fields provided")
	}

	sqlSet := "SET" + " " + strings.Join(sqlColumnValues, ",")
	sqlWhere := `WHERE id=$` + strconv.FormatInt(int64(len(values)+1), 10)
	sqlRet := `RETURNING id, name, key, description, created_by_id`

	sql := strings.Join([]string{"UPDATE projects", sqlSet, sqlWhere, sqlRet}, " ")

	values = append(values, projectId)

	row := q.Db.QueryRow(sql, values...)

	var project models.Project

	err := row.Scan(&project.Id, &project.Name, &project.Key, &project.Description, &project.Created_by_id)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (q *Queries) GetProjectDetails(projectId string) (*models.Project, error) {
	sql := `SELECT id, name, key, description, created_by_id, created_at FROM projects
	WHERE id=$1 
	LIMIT 1`

	var project models.Project

	row := q.Db.QueryRow(sql, projectId)

	err := row.Scan(&project.Id, &project.Name, &project.Key, &project.Description, &project.Created_by_id, &project.Created_at)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (q Queries) SelectProjectsForUser(userId string) ([]*ProjectDetails, error) {
	rows, err := q.Db.Query(`SELECT p.id, p.name, p.key, p.description, p.created_at, users.id, users.username FROM user_project_xref as xref
	INNER JOIN projects as p on p.id = xref.project_id
	LEFT JOIN users on users.id=p.created_by_id
	WHERE xref.user_id = $1`, userId)

	if err != nil {
		return nil, err
	}

	projects := []*ProjectDetails{}

	for rows.Next() {
		var p ProjectDetails
		err := rows.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.Created_at, &p.Creator.Id, &p.Creator.Username)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}

	return projects, nil
}

type JoinedProjectDTO struct {
	models.Project
	Created_by models.UserDTO `json:"created_by"`
}

func (q *Queries) GetJoinedProjectDetails(projectId string) (*ProjectDetails, error) {
	sql := `SELECT p.id, p.name, p.key, p.description, p.created_at, u.id, u.username from projects AS p
	INNER JOIN users AS u ON p.created_by_id=u.id
	WHERE p.id=$1
	`
	row := q.Db.QueryRow(sql, projectId)

	var p ProjectDetails
	u := &p.Creator

	err := row.Scan(&p.Id, &p.Name, &p.Key, &p.Description, &p.Created_at, &u.Id, &u.Username)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (q *Queries) GetProjectMembers(projectId string) ([]*models.User, error) {
	rows, err := q.Db.Query(`SELECT u.id, u.username, u.created_at FROM user_project_xref AS upxref
		INNER JOIN users as u ON u.id=upxref.user_id
	 	WHERE project_id=$1`, projectId)

	if err != nil {
		return nil, err
	}

	var users []*models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.Id, &user.Username, &user.Created_at)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (q *Queries) IsProjectMember(userId string, projectId string) (bool, error) {

	var count int
	row := q.Db.QueryRow(`SELECT COUNT(*) from user_project_xref WHERE user_id=$1 AND project_id=$2`, userId, projectId)

	err := row.Scan(&count)

	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}
