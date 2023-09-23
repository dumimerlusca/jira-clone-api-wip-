package queries

import (
	"encoding/json"
	"fmt"
	"jira-clone/packages/events"
	"jira-clone/packages/models"
)

func (q *Queries) RegisterEvent(eventSourceId string, data any) (*models.Event, error) {
	v, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	var event models.Event

	row := q.Db.QueryRow(`INSERT INTO events(source_id, data) VALUES($1, $2) RETURNING id, created_at, source_id, data`, eventSourceId, v)

	err = row.Scan(&event.Id, &event.Created_at, &event.Source_id, &event.Data)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &event, nil
}

type NewTicketUpdatedEventPayload struct {
	FieldName        string
	UpdatedById      string
	UpdateByUsername string
	FromValue        any
	FromDisplayValue string
	ToValue          any
	ToDisplayValue   string
	TicketId         string
}

func (q *Queries) NewTicketUpdatedEvent(payload NewTicketUpdatedEventPayload) (*models.Event, error) {
	event, err := q.RegisterEvent(events.SourceIdTicketUpdatedEvent, events.TicketUpdatedEventData{FieldName: payload.FieldName, TicketId: payload.TicketId, UpdatedBy: events.UserEventData{Username: payload.UpdateByUsername, Id: payload.UpdatedById}, From: &events.TicketUpdatedEventValues{Value: payload.FromValue, DisplayValue: payload.FromDisplayValue}, To: &events.TicketUpdatedEventValues{Value: payload.ToValue, DisplayValue: payload.ToDisplayValue}})
	if err != nil {
		return nil, err
	}

	return event, nil
}
