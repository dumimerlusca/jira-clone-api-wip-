package queries

import (
	"jira-clone/packages/random"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTicket(t *testing.T) {
	t.Run("should create a ticket", func(t *testing.T) {
		project, user := tQueries.CreateRandomProject(t)

		priority := 2
		story_points := 5
		desc := "Description"

		data := CreateTicketDTO{
			Title:         "TICKET 1",
			Type:          "bug",
			Project_id:    project.Id,
			Created_by_id: user.Id,
			Story_points:  &story_points,
			Description:   &desc,
			Priority:      &priority,
			Assignee_id:   &user.Id,
		}

		k, err := tQueries.CreateTicket(data)

		require.NoError(t, err)

		assert.NotEmpty(t, k.Number)
		assert.Equal(t, data.Title, k.Title)
		assert.Equal(t, data.Type, k.Type)
		assert.Equal(t, data.Project_id, k.Project_id)
		assert.Equal(t, data.Created_by_id, k.Created_by_id)
		assert.Equal(t, data.Description, k.Description)
		assert.Equal(t, data.Assignee_id, k.Assignee_id)
		assert.Equal(t, *data.Priority, k.Priority)
		assert.Equal(t, *data.Story_points, k.Story_points)
	})
}

func TestUpdateTicket(t *testing.T) {
	t.Run("should properly update a ticket", func(t *testing.T) {
		ticket := tQueries.CreateRandomTicket(t)

		user1, _ := tQueries.CreateRandomUser(t)

		title := random.RandomString(10)
		description := random.RandomString(10)
		story_points := int(random.RandomInt(0, 1000))
		priority := 3
		assigne_id := user1.Id
		status := "tested"
		ticketType := "epic"

		payload := UpdateTicketDTO{Title: &title, Type: &ticketType, Description: &description, Assignee_id: &assigne_id, Story_points: &story_points, Priority: &priority, Status: &status}

		updatedTicket, err := tQueries.UpdateTicket(ticket.Id, payload)

		require.NoError(t, err)

		assert.Equal(t, title, updatedTicket.Title)
		assert.Equal(t, ticketType, updatedTicket.Type)
		assert.Equal(t, description, *updatedTicket.Description)
		assert.Equal(t, story_points, updatedTicket.Story_points)
		assert.Equal(t, assigne_id, *updatedTicket.Assignee_id)
		assert.Equal(t, status, updatedTicket.Status)
		assert.NotEqual(t, ticket.Updated_at, updatedTicket.Updated_at)
	})
}

func TestFindTicketById(t *testing.T) {
	t.Run("should return the ticket if it exists", func(t *testing.T) {
		tkt := tQueries.CreateRandomTicket(t)
		ticket, err := tQueries.FindTicketById(tkt.Id)

		require.NoError(t, err)

		require.NotEmpty(t, ticket)
		require.Equal(t, tkt.Id, ticket.Id)

	})

	t.Run("should return error if ticket is not found", func(t *testing.T) {
		ticket, err := tQueries.FindTicketById(random.RandomString(30))

		require.Error(t, err)
		require.Empty(t, ticket)
	})
}
