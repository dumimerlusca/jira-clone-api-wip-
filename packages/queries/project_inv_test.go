package queries

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectSentProjectInvites(t *testing.T) {
	t.Run("should return an array with all invitations sent by a specific user", func(t *testing.T) {
		i, _ := tQueries.CreateRandomProjectInvite(t, "pending")

		sender_id := i.Sender.Id

		_, err := tQueries.CreateProjectInvitation(CreateProjectInvitationPayload{Receiver_id: i.Inv.Receiver_id, Project_id: i.Inv.Project_id, Sender_id: sender_id, Status: "pending"})

		require.NoError(t, err)

		inv, err := tQueries.SelectSentProjectInvites(sender_id)

		require.NoError(t, err)
		require.Equal(t, 2, len(inv))
	})
}

func TestReceivedProjectInvites(t *testing.T) {
	t.Run("should return an array with all invitations received", func(t *testing.T) {
		i, _ := tQueries.CreateRandomProjectInvite(t, "pending")

		receiverId := i.Receiver.Id

		_, err := tQueries.CreateProjectInvitation(CreateProjectInvitationPayload{Receiver_id: receiverId, Project_id: i.Inv.Project_id, Sender_id: i.Inv.Sender_id, Status: "pending"})

		require.NoError(t, err)

		inv, err := tQueries.SelectReceivedProjectInvites(receiverId)

		require.NoError(t, err)
		require.Equal(t, 2, len(inv))
	})
}
