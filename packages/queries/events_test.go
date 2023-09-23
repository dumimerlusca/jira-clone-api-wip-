package queries

import (
	"encoding/json"
	"jira-clone/packages/events"
	"jira-clone/packages/random"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterEvent(t *testing.T) {
	type eventData struct {
		Username string `json:"username"`
		Id       string `json:"id"`
	}

	data := eventData{Username: random.RandomString(10), Id: random.RandomString(10)}

	event, err := tQueries.RegisterEvent(events.SourceIdTicketUpdatedEvent, data)

	require.NoError(t, err)
	require.NotNil(t, event)

	var jData eventData

	err = json.Unmarshal([]byte(event.Data), &jData)

	require.NoError(t, err)
	assert.Equal(t, data.Username, jData.Username)
	assert.Equal(t, data.Id, jData.Id)
}

func TestNewTicketUpdatedEvent(t *testing.T) {
	data := NewTicketUpdatedEventPayload{FieldName: random.RandomString(10), UpdatedById: random.RandomString(10), UpdateByUsername: random.RandomString(10), FromValue: random.RandomString(10), FromDisplayValue: random.RandomString(10), ToValue: random.RandomString(10), ToDisplayValue: random.RandomString(10), TicketId: random.RandomString(10)}

	event, err := tQueries.NewTicketUpdatedEvent(data)

	require.NoError(t, err)

	assert.Equal(t, events.SourceIdTicketUpdatedEvent, event.Source_id)
	assert.NotEmpty(t, event.Created_at)
	assert.NotEmpty(t, event.Id)
	assert.NotEmpty(t, event.Data)

	var d events.TicketUpdatedEventData

	err = json.Unmarshal([]byte(event.Data), &d)

	require.NoError(t, err)

	assert.Equal(t, data.FieldName, d.FieldName)
	assert.Equal(t, data.FromValue, d.From.Value)
	assert.Equal(t, data.FromDisplayValue, d.From.DisplayValue)
	assert.Equal(t, data.ToValue, d.To.Value)
	assert.Equal(t, data.ToDisplayValue, d.To.DisplayValue)
	assert.Equal(t, data.UpdateByUsername, d.UpdatedBy.Username)
	assert.Equal(t, data.UpdatedById, d.UpdatedBy.Id)

}
