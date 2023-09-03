package queries

import "jira-clone/packages/models"

type CreateCommentDTO struct {
	Author_id string
	Ticket_id string
	Text      string
}

func (q *Queries) DeleteComment(id string) error {
	_, err := q.Db.Exec(`DELETE FROM comments WHERE id=$1`, id)

	return err
}

func (q *Queries) FindCommentById(id string) (*models.Comment, error) {
	row := q.Db.QueryRow(`SELECT id, ticket_id, author_id, text, created_at, updated_at FROM comments WHERE id=$1`, id)

	var c models.Comment

	err := row.Scan(&c.Id, &c.Ticket_id, &c.Author_id, &c.Text, &c.Created_at, &c.Updated_at)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (q *Queries) CreateComment(dto CreateCommentDTO) (*models.Comment, error) {
	var comment models.Comment

	row := q.Db.QueryRow(`INSERT INTO comments(author_id, ticket_id, text) VALUES($1, $2, $3) RETURNING id,author_id, ticket_id, text, created_at, updated_at`, dto.Author_id, dto.Ticket_id, dto.Text)

	err := row.Scan(&comment.Id, &comment.Author_id, &comment.Ticket_id, &comment.Text, &comment.Created_at, &comment.Updated_at)

	if err != nil {
		return nil, err
	}

	return &comment, nil

}

type CommentAuthor struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type JoinedComment struct {
	Id         string        `json:"id"`
	Ticket_id  string        `json:"ticket_id"`
	Text       string        `json:"text"`
	Created_at string        `json:"created_at"`
	Updated_at string        `json:"updated_at"`
	Author     CommentAuthor `json:"author"`
}

func (q *Queries) SelectJoinedTicketComments(ticketId string) ([]JoinedComment, error) {
	rows, err := q.Db.Query(`SELECT c.id, c.text,c.ticket_id, c.created_at, c.updated_at,u.id, u.username FROM comments AS c
		INNER JOIN users AS u ON u.id=c.author_id
		WHERE ticket_id=$1`, ticketId)

	if err != nil {
		return nil, err
	}

	comments := []JoinedComment{}

	for rows.Next() {
		var comment JoinedComment

		err := rows.Scan(&comment.Id, &comment.Text, &comment.Ticket_id, &comment.Created_at, &comment.Updated_at, &comment.Author.Id, &comment.Author.Username)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
