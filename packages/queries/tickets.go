package queries

import (
	"fmt"
	"jira-clone/packages/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateTicketDTO struct {
	Project_id    string
	Created_by_id string
	Title         string
	Type          string
	Assignee_id   *string
	Component_id  *string
	Story_points  *int
	Description   *string
	Priority      *int
}

func (d *CreateTicketDTO) Validate() error {
	if d.Title == "" {
		return fmt.Errorf("title required")
	}

	if d.Type == "" {
		return fmt.Errorf("type required")
	}

	return nil
}

func (q *Queries) CreateTicket(d CreateTicketDTO) (*models.Ticket, error) {
	id := uuid.NewString()

	fields := "id,title, project_id, created_by_id, updated_by_id, type"
	values := []any{id, d.Title, d.Project_id, d.Created_by_id, d.Created_by_id, d.Type}

	handleField := func(name string, value any) {
		values = append(values, value)
		fields = fields + "," + name
	}

	if d.Assignee_id != nil {
		handleField("assignee_id", &d.Assignee_id)
	}

	if d.Component_id != nil {
		handleField("component_id", &d.Component_id)
	}

	if d.Story_points != nil {
		handleField("story_points", &d.Story_points)
	}

	if d.Description != nil {
		handleField("description", &d.Description)
	}

	if d.Priority != nil {
		handleField("priority", &d.Priority)
	}

	sql := `INSERT INTO tickets(` + fields + ") "
	v := []string{}

	for i := range values {
		v = append(v, "$"+strconv.FormatInt(int64(i+1), 10))
	}

	sqlValues := `VALUES(` + strings.Join(v, ",") + ") "

	sql = sql + sqlValues + `RETURNING id, number,type,priority, title, story_points, description, status, created_by_id, assignee_id, project_id, component_id, updated_by_id, created_at, updated_at`

	row := q.Db.QueryRow(sql, values...)

	var t models.Ticket

	err := row.Scan(&t.Id, &t.Number, &t.Type, &t.Priority, &t.Title, &t.Story_points, &t.Description, &t.Status, &t.Created_by_id, &t.Assignee_id, &t.Project_id, &t.Component_id, &t.Updated_by_id, &t.Created_at, &t.Updated_at)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

type UpdateTicketDTO struct {
	Title         *string
	Type          *string
	Assignee_id   *string
	Component_id  *string
	Story_points  *int
	Description   *string
	Priority      *int
	Status        *string
	Updated_by_id *string
}

func (q *Queries) UpdateTicket(ticketId string, d UpdateTicketDTO) (*models.Ticket, error) {

	values := []any{time.Now()}
	sqlColumnValues := []string{"updated_at=$1"}

	handleField := func(name string, value any) {
		values = append(values, value)
		s := name + "=$" + strconv.FormatInt(int64(len(values)), 10)
		sqlColumnValues = append(sqlColumnValues, s)
	}

	if d.Title != nil {
		handleField("title", *d.Title)
	}

	if d.Assignee_id != nil {
		handleField("assignee_id", *d.Assignee_id)
	}

	if d.Component_id != nil {
		handleField("component_id", *d.Component_id)
	}

	if d.Description != nil {
		handleField("description", *d.Description)
	}
	if d.Priority != nil {
		handleField("priority", *d.Priority)
	}

	if d.Story_points != nil {
		handleField("story_points", *d.Story_points)
	}
	if d.Status != nil {
		handleField("status", *d.Status)
	}
	if d.Updated_by_id != nil {
		handleField("updated_by_id", *d.Updated_by_id)
	}
	if d.Type != nil {
		handleField("type", *d.Type)
	}

	set := strings.Join(sqlColumnValues, ",")
	where := ` WHERE id=$` + strconv.FormatInt(int64((len(values)+1)), 10)
	sql := `UPDATE tickets SET ` + set + where + " RETURNING id, type,priority, title, story_points, description, status, created_by_id, assignee_id, project_id, component_id, updated_by_id, created_at, updated_at"

	values = append(values, ticketId)

	fmt.Print(sql)

	row := q.Db.QueryRow(sql, values...)

	var t models.Ticket

	err := row.Scan(&t.Id, &t.Type, &t.Priority, &t.Title, &t.Story_points, &t.Description, &t.Status, &t.Created_by_id, &t.Assignee_id, &t.Project_id, &t.Component_id, &t.Updated_by_id, &t.Created_at, &t.Updated_at)

	if err != nil {
		return nil, err
	}

	return &t, nil

}

func (q *Queries) FindTicketById(id string) (*models.Ticket, error) {
	var t models.Ticket

	row := q.Db.QueryRow(`SELECT id,priority, title, story_points, description, status, created_by_id, assignee_id, project_id, component_id, updated_by_id, created_at, updated_at FROM tickets WHERE id=$1`, id)

	err := row.Scan(&t.Id, &t.Priority, &t.Title, &t.Story_points, &t.Description, &t.Status, &t.Created_by_id, &t.Assignee_id, &t.Project_id, &t.Component_id, &t.Updated_by_id, &t.Created_at, &t.Updated_at)

	if err != nil {
		return nil, err
	}

	return &t, nil

}

func (q *Queries) FindTicketKeyById(id string) (*string, error) {
	row := q.Db.QueryRow(`SELECT key FROM tickets_view WHERE id=$1`, id)

	var key string
	err := row.Scan(&key)

	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (q *Queries) GetProjectIdsWhereUserIsMember(userId string) ([]string, error) {
	rows, err := q.Db.Query(`SELECT project_id FROM user_project_xref WHERE user_id=$1`, userId)

	if err != nil {
		return nil, err
	}

	projecIds := []string{}

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		projecIds = append(projecIds, id)
	}

	return projecIds, nil
}

func (q *Queries) GetTicketDetailsByKeyForUser(ticketKey string, userId string) (*TicketDetails, error) {
	projectIds, err := q.GetProjectIdsWhereUserIsMember(userId)

	if err != nil {
		return nil, err
	}

	if len(projectIds) == 0 {
		return nil, fmt.Errorf("user is not part of any project")
	}

	mappedIds := []string{}
	for _, str := range projectIds {
		mappedIds = append(mappedIds, "'"+str+"'")
	}

	sql := `SELECT 
			id,
			key,
			type,
			priority,
			title,
			story_points,
			description,
			status,
			component_id,
			created_at,
			updated_at,
			creator_id,
			creator_username,
			assignee_id,
			assignee_username
	 FROM tickets_view WHERE key=$1 AND project_id IN ` + fmt.Sprintf("(%v)", strings.Join(mappedIds, ","))

	row := q.Db.QueryRow(sql, ticketKey)

	var createdBy UserItem
	var assignee UserItem

	t := TicketDetails{Creator: &createdBy, Assignee: &assignee}

	err = row.Scan(&t.Id, &t.Key, &t.Type, &t.Priority, &t.Title, &t.Story_points, &t.Description, &t.Status, &t.Component_id, &t.Created_at, &t.Updated_at, &t.Creator.Id, &t.Creator.Username, &t.Assignee.Id, &t.Assignee.Username)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
