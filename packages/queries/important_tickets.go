package queries

import "jira-clone/packages/models"

func (q *Queries) FindImportantTicket(userId string, ticketId string) (*models.ImportantTicket, error) {
	row := q.Db.QueryRow(`SELECT ticket_id, user_id, project_id FROM important_tickets WHERE user_id=$1 AND ticket_id=$2`, userId, ticketId)

	var t models.ImportantTicket

	err := row.Scan(&t.Ticket_id, &t.User_id, &t.Project_id)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (q *Queries) CreateImportantTicket(userId string, ticketId string, projectId string) (*models.ImportantTicket, error) {
	row := q.Db.QueryRow(`INSERT INTO important_tickets(ticket_id, user_id, project_id) VALUES($1, $2, $3) RETURNING ticket_id, user_id, project_id`, ticketId, userId, projectId)

	var t models.ImportantTicket

	err := row.Scan(&t.Ticket_id, &t.User_id, &t.Project_id)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (q *Queries) DeleteImportantTicket(userId string, ticketId string) error {
	_, err := q.Db.Exec(`DELETE FROM important_tickets WHERE user_id=$1 AND ticket_id=$2`, userId, ticketId)

	return err
}
