package queries

import "jira-clone/packages/models"

func (q *Queries) GetUsers() ([]models.User, error) {
	db := q.Db

	sql := `SELECT id, username, created_at FROM users`

	rows, err := db.Query(sql)

	if err != nil {
		return nil, err
	}

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Username, &user.Created_at)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (q *Queries) CreateUser(id string, username string, password string) (*models.User, error) {
	db := q.Db

	sql := `INSERT INTO users(id, username, password)
	VALUES($1, $2, $3) RETURNING id, username, created_at`

	rows, err := db.Query(sql, id, username, password)

	if err != nil {
		return nil, err
	}

	var user models.User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Created_at)
		if err != nil {
			return nil, err
		}
	}

	return &user, err
}

func (q *Queries) FindUserByUsername(username string, includePassword bool) (*models.User, error) {
	db := q.Db

	var sql string
	var args []any
	var user models.User

	if includePassword {
		sql = `SELECT id, username, created_at, password FROM users WHERE username = $1`
		args = append(args, &user.Id, &user.Username, &user.Created_at, &user.Password)
	} else {
		sql = `SELECT id, username, created_at FROM users WHERE username = $1`
		args = append(args, &user.Id, &user.Username, &user.Created_at)
	}

	row := db.QueryRow(sql, username)

	err := row.Scan(args...)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
