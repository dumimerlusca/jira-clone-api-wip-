package queries

import (
	"fmt"
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

func TestFindTicketKey(t *testing.T) {
	t.Run("should return err", func(t *testing.T) {
		key, err := tQueries.FindTicketKeyById("non exitent ticket")

		assert.Error(t, err)
		assert.Empty(t, key)
	})

	t.Run("should return the key", func(t *testing.T) {
		p, _ := tQueries.CreateRandomProject(t)
		ticket := tQueries.CreateRandomTicketForProject(t, p.Id, p.Created_by_id)
		key, err := tQueries.FindTicketKeyById(ticket.Id)

		require.NoError(t, err)
		require.NotEmpty(t, key)

		assert.Equal(t, fmt.Sprintf("%v-1", p.Key), *key)
	})
}
func TestGetProjectIdsWhereUserIsMember(t *testing.T) {

	t.Run("should return list with ids", func(t *testing.T) {
		u, _ := tQueries.CreateRandomUser(t)
		u2, _ := tQueries.CreateRandomUser(t)

		tQueries.CreateRandomProjectForUser(t, u.Id)
		tQueries.CreateRandomProjectForUser(t, u.Id)
		tQueries.CreateRandomProjectForUser(t, u.Id)
		tQueries.CreateRandomProjectForUser(t, u2.Id)

		ids, err := tQueries.GetProjectIdsWhereUserIsMember(u.Id)

		for _, val := range ids {
			assert.NotEmpty(t, val)
		}

		require.NoError(t, err)

		require.Equal(t, 3, len(ids))
	})
}

func TestGetTicketDetailsByKeyForUser(t *testing.T) {
	t.Run("should work", func(t *testing.T) {
		ticket := tQueries.CreateRandomTicket(t)
		key, err := tQueries.FindTicketKeyById(ticket.Id)

		require.NoError(t, err)

		details, err := tQueries.GetTicketDetailsByKeyForUser(*key, ticket.Created_by_id)
		require.NoError(t, err)
		require.NotNil(t, details)

		assert.Equal(t, *key, details.Key)
		assert.Equal(t, ticket.Id, details.Id)
		assert.Equal(t, ticket.Title, details.Title)
	})
}
