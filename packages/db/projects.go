package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type CreateProjectDTO struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Leader_id   string `json:"leader_id"`
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

func (q *Queries) CreateProject(data CreateProjectDTO) (*Project, error) {
	id := uuid.NewString()

	sql := `INSERT INTO projects(id, name, key, description, leader_id) VALUES($1, $2, $3, $4, $5) RETURNING id, name, key, description, leader_id, created_at`

	row := q.Db.QueryRow(sql, id, data.Name, data.Key, data.Description, data.Leader_id)

	var project Project

	err := row.Scan(&project.Id, &project.Name, &project.Key, &project.Description, &project.Leader_id, &project.Created_at)

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

func (q *Queries) UpdateProject(projectId string, data UpdateProjectDTO) (*Project, error) {
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
	sqlRet := `RETURNING id, name, key, description, leader_id`

	sql := strings.Join([]string{"UPDATE projects", sqlSet, sqlWhere, sqlRet}, " ")

	values = append(values, projectId)

	row := q.Db.QueryRow(sql, values...)

	var project Project

	err := row.Scan(&project.Id, &project.Name, &project.Key, &project.Description, &project.Leader_id)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (q *Queries) GetProjectDetails(projectId string) (*Project, error) {
	sql := `SELECT id, name, key, description, leader_id, created_at FROM projects
	WHERE id=$1 
	LIMIT 1`

	var project Project

	row := q.Db.QueryRow(sql, projectId)

	err := row.Scan(&project.Id, &project.Name, &project.Key, &project.Description, &project.Leader_id, &project.Created_at)

	if err != nil {
		return nil, err
	}

	return &project, nil
}
