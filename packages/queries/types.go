package queries

type TicketDetails struct {
	Id           string  `json:"id"`
	Project_id   string  `json:"project_id"`
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
	Important    bool    `json:"important"`

	Creator  *UserItem `json:"creator"`
	Assignee *UserItem `json:"assignee"`
}

type ProjectDetails struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Description string   `json:"description"`
	Creator     UserItem `json:"creator"`
	Created_at  string   `json:"created_at"`
}

type UserItem struct {
	Id       *string `json:"id"`
	Username *string `json:"username"`
}
