package queries

type TicketDetails struct {
	Id           string  `json:"id"`
	Key          string  `json:"key"`
	Type         string  `json:"type"`
	Priority     int     `json:"priority"`
	Title        string  `json:"title"`
	Story_points int     `json:"story_points"`
	Description  *string `json:"description"`
	Status       string  `json:"status"`
	Component_id *string `json:"component_id"`
	Created_at   string  `json:"created_at"`
	Updated_at   string  `json:"updated_at"`

	Creator  *UserItem `json:"creator"`
	Assignee *UserItem `json:"assignee"`
}

type UserItem struct {
	Id       *string `json:"id"`
	Username *string `json:"username"`
}
