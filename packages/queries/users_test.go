package queries

import (
	"jira-clone/packages/random"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	id := random.RandomString(10)
	username := random.RandomString(10)
	password := random.RandomString(10)
	user, err := tQueries.CreateUser(id, username, password)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, id, user.Id)
	require.Equal(t, username, user.Username)
	// Password should not be returned back
	require.Zero(t, user.Password)
	require.NotZero(t, user.Created_at)

	row := tQueries.Db.QueryRow("SELECT password from users WHERE id=$1", id)

	var pswd string

	err = row.Scan(&pswd)

	require.NoError(t, err)

	require.Equal(t, password, pswd)
}

func TestFindUserByUsername(t *testing.T) {
	u, _ := tQueries.CreateRandomUser(t)

	// Password included
	user, err := tQueries.FindUserByUsername(u.Username, true)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, u.Username)
	require.NotZero(t, user.Password)

	// Password not included
	user, err = tQueries.FindUserByUsername(u.Username, false)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, u.Username)
	require.Zero(t, user.Password)

}
