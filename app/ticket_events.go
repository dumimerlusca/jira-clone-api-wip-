package app

import (
	"fmt"
	"jira-clone/packages/models"
	"jira-clone/packages/queries"
	"strconv"
)

func (app *application) registerTicketUpdatedEvents(old *models.Ticket, new *models.Ticket) {
	user, err := app.queries.FindUserById(new.Updated_by_id, false)

	if err != nil {
		fmt.Println(err)
		return
	}

	if old.Description == nil && new.Description != nil {
		app.queries.NewTicketUpdatedEvent(queries.NewTicketUpdatedEventPayload{FieldName: "Description", UpdatedById: user.Id, UpdateByUsername: user.Username, TicketId: new.Id, ToValue: *new.Description, ToDisplayValue: *new.Description})
	}

	if old.Description != nil && new.Description == nil {
		app.queries.NewTicketUpdatedEvent(queries.NewTicketUpdatedEventPayload{FieldName: "Description", UpdatedById: user.Id, UpdateByUsername: user.Username, TicketId: new.Id, FromValue: *old.Description, FromDisplayValue: *old.Description})
	}

	if old.Description != nil && new.Description != nil && *old.Description != *new.Description {
		app.queries.NewTicketUpdatedEvent(queries.NewTicketUpdatedEventPayload{FieldName: "Description", UpdatedById: user.Id, UpdateByUsername: user.Username, TicketId: new.Id, FromValue: *old.Description, FromDisplayValue: *old.Description, ToValue: *new.Description, ToDisplayValue: *new.Description})
	}

	if old.Title != new.Title {
		app.queries.NewTicketUpdatedEvent(queries.NewTicketUpdatedEventPayload{FieldName: "Title", UpdatedById: user.Id, UpdateByUsername: user.Username, TicketId: new.Id, FromValue: old.Title, FromDisplayValue: old.Title, ToValue: new.Title, ToDisplayValue: new.Title})
	}
	if old.Story_points != new.Story_points {
		app.queries.NewTicketUpdatedEvent(queries.NewTicketUpdatedEventPayload{FieldName: "Story Points", UpdatedById: user.Id, UpdateByUsername: user.Username, TicketId: new.Id, FromValue: old.Story_points, FromDisplayValue: fmt.Sprintf(`%v Story Points`, old.Story_points), ToValue: new.Story_points, ToDisplayValue: fmt.Sprintf(`%v Story Points`, new.Story_points)})
	}

	if old.Status != new.Status {
		app.queries.NewTicketUpdatedEvent(queries.NewTicketUpdatedEventPayload{FieldName: "Status", UpdatedById: user.Id, UpdateByUsername: user.Username, TicketId: new.Id, FromValue: old.Status, FromDisplayValue: old.Status, ToValue: new.Status, ToDisplayValue: new.Status})
	}

	if old.Type != new.Type {
		app.queries.NewTicketUpdatedEvent(queries.NewTicketUpdatedEventPayload{FieldName: "Ticket type", UpdatedById: user.Id, UpdateByUsername: user.Username, TicketId: new.Id, FromValue: old.Type, FromDisplayValue: old.Type, ToValue: new.Type, ToDisplayValue: new.Type})
	}

	if old.Priority != new.Priority {
		app.queries.NewTicketUpdatedEvent(queries.NewTicketUpdatedEventPayload{FieldName: "Priority", UpdatedById: user.Id, UpdateByUsername: user.Username, TicketId: new.Id, FromValue: old.Priority, FromDisplayValue: getTicketPriorityLabel(old.Priority), ToValue: new.Priority, ToDisplayValue: getTicketPriorityLabel(new.Priority)})
	}
}

func getTicketPriorityLabel(priority int) string {
	switch priority {
	case 0:
		return "Highest"
	case 1:
		return "High"
	case 2:
		return "Medium"
	case 3:
		return "Low"
	case 4:
		return "Lowest"
	default:
		return strconv.FormatInt(int64(priority), 10)
	}

}
