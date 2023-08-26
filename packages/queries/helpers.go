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

func (q *Queries) CreateRandomProject(t *testing.T) (*models.Project, *models.User) {
	name := random.RandomString(20)
	description := random.RandomString(20)
	key := random.RandomString(4)
	user, _ := q.CreateRandomUser(t)

	data := CreateProjectDTO{Name: name, Description: description, Key: key, Created_by_id: user.Id}
	project, err := q.CreateProject(data)

	require.NoError(t, err)

	return project, user
}
