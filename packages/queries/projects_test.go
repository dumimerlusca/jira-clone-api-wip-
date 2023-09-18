package queries

import (
	"jira-clone/packages/random"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestSelectProjectsForUser(t *testing.T) {
	t.Run("should return an array with all the projects a user is part of", func(t *testing.T) {
		p, user := tQueries.CreateRandomProject(t)

		projects, err := tQueries.SelectProjectsForUser(user.Id)

		require.NoError(t, err)
		require.Equal(t, 1, len(projects))

		item := projects[0]

		assert.Equal(t, p.Id, item.Id)
		assert.Equal(t, p.Key, item.Key)
		assert.Equal(t, p.Description, item.Description)
		assert.Equal(t, p.Name, item.Name)
		assert.Equal(t, p.Created_at, item.Created_at)
		assert.Equal(t, p.Created_by_id, item.Created_by_id)
	})

	t.Run("should return empty slice if no results are found", func(t *testing.T) {
		projects, err := tQueries.SelectProjectsForUser("non existent user")
		require.NoError(t, err)
		require.Equal(t, 0, len(projects))
	})
}
