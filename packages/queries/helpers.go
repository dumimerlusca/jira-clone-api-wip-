package queries

import (
	"jira-clone/packages/models"
	"jira-clone/packages/random"
	"testing"

	"github.com/stretchr/testify/require"
)

func (q *Queries) CreateRandomUser(t *testing.T) (user *models.User, decodedPassword string) {
	id := random.RandomString(20)
	username := random.RandomString(20)
	password := random.RandomString(6)
	user, err := q.CreateUser(id, username, password)

	require.NoError(t, err)

	return user, password
}

func (q *Queries) CreateRandomProjectForUser(t *testing.T, userId string) *models.Project {
	name := random.RandomString(20)
	description := random.RandomString(20)
	key := random.RandomString(4)

	data := CreateProjectDTO{Name: name, Description: description, Key: key, Created_by_id: userId}
	project, err := q.CreateProject(data)
	q.CreateUserProjectXref(userId, project.Id)

	require.NoError(t, err)

	return project
}

func (q *Queries) CreateRandomProject(t *testing.T) (*models.Project, *models.User) {
	name := random.RandomString(20)
	description := random.RandomString(20)
	key := random.RandomString(4)
	user, _ := q.CreateRandomUser(t)

	data := CreateProjectDTO{Name: name, Description: description, Key: key, Created_by_id: user.Id}
	project, err := q.CreateProject(data)
	q.CreateUserProjectXref(user.Id, project.Id)

	require.NoError(t, err)

	return project, user
}

type CreateRandomProjectInviteReturnValue struct {
	Inv      *models.ProjectInvitation
	Project  *models.Project
	Sender   *models.User
	Receiver *models.User
}

func (q *Queries) CreateRandomProjectInvite(t *testing.T, status string) (*CreateRandomProjectInviteReturnValue, error) {
	project, sender := q.CreateRandomProject(t)
	receiver, _ := q.CreateRandomUser(t)

	p := CreateProjectInvitationPayload{Status: status, Receiver_id: receiver.Id, Project_id: project.Id, Sender_id: sender.Id}

	inv, err := q.CreateProjectInvitation(p)

	return &CreateRandomProjectInviteReturnValue{Inv: inv, Project: project, Receiver: receiver, Sender: sender}, err
}

func (q *Queries) CreateRandomTicketForProject(t *testing.T, projectId string, userId string) *models.Ticket {
	data := CreateTicketDTO{Project_id: projectId, Type: "bug", Created_by_id: userId, Title: random.RandomString(20)}
	ticket, err := q.CreateTicket(data)

	require.NoError(t, err)

	return ticket
}

func (q *Queries) CreateRandomTicket(t *testing.T) *models.Ticket {
	project, user := q.CreateRandomProject(t)
	data := CreateTicketDTO{Project_id: project.Id, Type: "bug", Created_by_id: user.Id, Title: random.RandomString(20)}
	ticket, err := q.CreateTicket(data)

	require.NoError(t, err)

	return ticket
}

func (q *Queries) CreateRandomComment(t *testing.T) *models.Comment {
	ticket := q.CreateRandomTicket(t)
	d := CreateCommentDTO{Author_id: ticket.Created_by_id, Text: random.RandomString(20), Ticket_id: ticket.Id}
	com, err := q.CreateComment(d)

	require.NoError(t, err)

	return com
}
