package db

func (q *Queries) GetUsers() ([]User, error) {
	db := q.Db

	sql := `SELECT id, username, created_at FROM users`

	rows, err := db.Query(sql)

	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Created_at)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (q *Queries) CreateUser(id string, username string, password string) (*User, error) {
	db := q.Db

	sql := `INSERT INTO users(id, username, password)
	VALUES($1, $2, $3) RETURNING id, username, created_at`

	rows, err := db.Query(sql, id, username, password)

	if err != nil {
		return nil, err
	}

	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Created_at)
		if err != nil {
			return nil, err
		}
	}

	return &user, err
}

func (q *Queries) FindUserByUsername(username string, includePassword bool) (*User, error) {
	db := q.Db

	var sql string

	if includePassword {
		sql = `SELECT id, username, created_at, password FROM users WHERE username = $1`
	} else {
		sql = `SELECT id, username, created_at FROM users WHERE username = $1`
	}

	row := db.QueryRow(sql, username)

	var user User

	err := row.Scan(&user.Id, &user.Username, &user.Created_at, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
