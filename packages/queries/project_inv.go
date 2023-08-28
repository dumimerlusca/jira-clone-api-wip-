package queries

import (
	"fmt"
	"jira-clone/packages/models"
)

type CreateProjectInvitationPayload struct {
	Receiver_id string
	Project_id  string
	Sender_id   string
	Status      string
}

func (q *Queries) CreateProjectInvitation(p CreateProjectInvitationPayload) (*models.ProjectInvitation, error) {
	var inv models.ProjectInvitation

	row := q.Db.QueryRow(`INSERT INTO project_invitations(receiver_id, project_id, sender_id, status)
	VALUES($1, $2, $3, $4) RETURNING id, receiver_id, project_id, sender_id, status`, p.Receiver_id, p.Project_id, p.Sender_id, p.Status)

	err := row.Scan(&inv.Id, &inv.Receiver_id, &inv.Project_id, &inv.Sender_id, &inv.Status)

	if err != nil {
		return nil, err
	}

	return &inv, err
}

func (q *Queries) FindProjectInvitationById(id string) (*models.ProjectInvitation, error) {
	var inv models.ProjectInvitation

	row := q.Db.QueryRow(`SELECT id, receiver_id, project_id, sender_id, status FROM project_invitations WHERE id=$1 `, id)

	err := row.Scan(&inv.Id, &inv.Receiver_id, &inv.Project_id, &inv.Sender_id, &inv.Status)

	if err != nil {
		return nil, err
	}

	return &inv, err

}

func (q *Queries) SelectPendingProjectInvitationsCount(receiver_id string, project_id string) (*int, error) {
	var count int
	row := q.Db.QueryRow(`SELECT COUNT(*) FROM project_invitations WHERE project_id=$1 AND receiver_id=$2 AND status='pending'`, project_id, receiver_id)
	err := row.Scan(&count)

	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (q *Queries) UpdateProjectInvitationStatus(id string, status string) error {
	result, err := q.Db.Exec(`UPDATE project_invitations SET status=$1 WHERE id=$2`, status, id)

	if err != nil {
		return err
	}

	rowCount, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowCount == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

type ProjectInvitationJoined struct {
	Id       string `json:"id"`
	Status   string `json:"status"`
	Receiver struct {
		Id       string `json:"id,omitempty"`
		Username string `json:"username,omitempty"`
	} `json:"receiver"`
	Sender struct {
		Id       string `json:"id,omitempty"`
		Username string `json:"username,omitempty"`
	} `json:"sender"`
	Project struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Key  string `json:"key"`
	} `json:"project"`
}

func (q *Queries) SelectSentProjectInvites(userId string) ([]ProjectInvitationJoined, error) {
	sql := `SELECT 
		p_inv.id,
		p_inv.status,
		u.id,
		u.username,
		p.id,
		p.name,
		p.key
	FROM project_invitations AS p_inv
	INNER JOIN users as u ON u.id=p_inv.receiver_id
	INNER JOIN projects as p ON p.id=p_inv.project_id
	WHERE sender_id=$1`

	rows, err := q.Db.Query(sql, userId)

	if err != nil {
		return nil, err
	}

	var invitations []ProjectInvitationJoined

	for rows.Next() {
		var inv ProjectInvitationJoined

		err := rows.Scan(&inv.Id, &inv.Status, &inv.Receiver.Id, &inv.Project.Id, &inv.Receiver.Username, &inv.Project.Name, &inv.Project.Key)

		if err != nil {
			return nil, err
		}

		invitations = append(invitations, inv)
	}

	return invitations, nil
}

func (q *Queries) SelectReceivedProjectInvites(userId string) ([]ProjectInvitationJoined, error) {
	sql := `SELECT 
		p_inv.id,
		p_inv.status,
		u.id,
		u.username,
		p.id,
		p.name,
		p.key
	FROM project_invitations AS p_inv
	INNER JOIN users as u ON u.id=p_inv.sender_id
	INNER JOIN projects as p ON p.id=p_inv.project_id
	WHERE receiver_id=$1`

	rows, err := q.Db.Query(sql, userId)

	if err != nil {
		return nil, err
	}

	var invitations []ProjectInvitationJoined

	for rows.Next() {
		var inv ProjectInvitationJoined

		err := rows.Scan(&inv.Id, &inv.Status, &inv.Sender.Id, &inv.Sender.Username, &inv.Project.Id, &inv.Project.Name, &inv.Project.Key)

		if err != nil {
			return nil, err
		}

		invitations = append(invitations, inv)
	}

	return invitations, nil
}
