package queries

import (
	"jira-clone/packages/random"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateComment(t *testing.T) {
	t.Run("should save comment to database and return it back", func(t *testing.T) {
		ticket := tQueries.CreateRandomTicket(t)

		dto := CreateCommentDTO{Text: random.RandomString(30), Author_id: ticket.Created_by_id, Ticket_id: ticket.Id}

		comment, err := tQueries.CreateComment(dto)

		require.NoError(t, err)

		assert.Equal(t, dto.Text, comment.Text)
		assert.Equal(t, dto.Author_id, comment.Author_id)
		assert.Equal(t, dto.Ticket_id, comment.Ticket_id)
		assert.NotZero(t, comment.Created_at)
		assert.NotZero(t, comment.Updated_at)

	})
}

func TestSelectJoinedTicketComments(t *testing.T) {
	t.Run("should return an array of comment", func(t *testing.T) {
		ticket := tQueries.CreateRandomTicket(t)

		d := CreateCommentDTO{Author_id: ticket.Created_by_id, Ticket_id: ticket.Id, Text: "asdasdads"}
		tQueries.CreateComment(d)
		tQueries.CreateComment(d)
		tQueries.CreateComment(d)

		comments, err := tQueries.SelectJoinedTicketComments(ticket.Id)

		require.NoError(t, err)

		require.Equal(t, 3, len(comments))

		c := comments[0]

		assert.Equal(t, d.Text, c.Text)
	})
}

func TestFindCommentById(t *testing.T) {
	t.Run("should return the comment", func(t *testing.T) {
		ticket := tQueries.CreateRandomTicket(t)
		d := CreateCommentDTO{Text: "sadasd", Author_id: ticket.Created_by_id, Ticket_id: ticket.Id}
		c, _ := tQueries.CreateComment(d)

		comment, err := tQueries.FindCommentById(c.Id)

		require.NoError(t, err)

		require.Equal(t, c.Id, comment.Id)

	})
}

func TestUpdateComment(t *testing.T) {
	t.Run("should update comment and return it", func(t *testing.T) {
		com := tQueries.CreateRandomComment(t)

		payload := UpdateCommentPayload{Text: random.RandomString(20)}

		updatedCom, err := tQueries.UpdateComment(com.Id, payload)

		require.NoError(t, err)

		require.Equal(t, payload.Text, updatedCom.Text)
		require.NotEqual(t, com.Updated_at, updatedCom.Updated_at)
	})
}
