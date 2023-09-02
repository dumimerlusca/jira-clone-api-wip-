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

func TestIsProjectMember(t *testing.T) {
	t.Run("should return true if user is project member", func(t *testing.T) {
		project, _ := tQueries.CreateRandomProject(t)
		user, _ := tQueries.CreateRandomUser(t)
		tQueries.Db.Exec(`INSERT INTO user_project_xref(user_id, project_id) VALUES($1, $2)`, user.Id, project.Id)
		v, err := tQueries.IsProjectMember(user.Id, project.Id)
		require.NoError(t, err)
		require.Equal(t, true, v)
	})
	t.Run("should return false if user is not project member", func(t *testing.T) {
		project, _ := tQueries.CreateRandomProject(t)
		user, _ := tQueries.CreateRandomUser(t)
		v, err := tQueries.IsProjectMember(user.Id, project.Id)
		require.NoError(t, err)
		require.Equal(t, false, v)
	})
}
