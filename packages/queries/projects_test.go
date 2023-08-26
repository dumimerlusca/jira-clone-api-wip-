package queries

import (
	"jira-clone/packages/random"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateProject(t *testing.T) {
	user, _ := tQueries.CreateRandomUser(t)

	args := CreateProjectDTO{Name: random.RandomString(10), Key: random.RandomString(4), Description: random.RandomString(20), Created_by_id: user.Id}

	project, err := tQueries.CreateProject(args)

	require.NoError(t, err)
	require.NotEmpty(t, project)
	require.Equal(t, args.Name, project.Name)
	require.Equal(t, args.Key, project.Key)
	require.Equal(t, args.Description, project.Description)
	require.Equal(t, args.Created_by_id, user.Id)

	require.NotZero(t, project.Id)
	require.NotZero(t, project.Created_at)
}
