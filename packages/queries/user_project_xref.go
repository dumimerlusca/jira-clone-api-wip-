package queries

import "jira-clone/packages/models"

func (q *Queries) CreateUserProjectXref(userId string, projectId string) (*models.UserProjectXref, error) {
	row := q.Db.QueryRow(`INSERT INTO user_project_xref(user_id, project_id) VALUES($1, $2) RETURNING user_Id, project_id`, userId, projectId)

	var xRef models.UserProjectXref

	err := row.Scan(&xRef.User_id, &xRef.Project_id)

	if err != nil {
		return nil, err
	}

	return &xRef, nil
}

func (q *Queries) IsUserInProject(userId string, projectId string) (int, error) {
	row := q.Db.QueryRow(`SELECT COUNT(*) FROM user_project_xref WHERE user_id=$1 AND project_id=$2`, userId, projectId)

	var count int

	err := row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
