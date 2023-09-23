package events

const SourceIdTicketUpdatedEvent = "1"

type UserEventData struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}

type TicketUpdatedEventValues struct {
	Value        any    `json:"value"`
	DisplayValue string `json:"displayValue"`
}

type TicketUpdatedEventData struct {
	FieldName string                    `json:"fieldName"`
	From      *TicketUpdatedEventValues `json:"from"`
	To        *TicketUpdatedEventValues `json:"to"`
	UpdatedBy UserEventData             `json:"updatedBy"`
	TicketId  string                    `json:"ticketId"`
}
