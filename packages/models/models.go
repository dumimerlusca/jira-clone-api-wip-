package models

type UserDTO struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Created_at string `json:"created_at"`
}

type User struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Created_at string `json:"created_at"`
	Password   string `json:"password,omitempty"`
}

type Project struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Key           string `json:"key"`
	Description   string `json:"description"`
	Created_by_id string `json:"created_by_id"`
	Created_at    string `json:"created_at"`
}

type ProjectInvitation struct {
	Id          string `json:"id"`
	Receiver_id string `json:"receiver_id"`
	Project_id  string `json:"project_id"`
	Sender_id   string `json:"sender_id"`
	Status      string `json:"status"`
	Created_at  string `json:"created_at"`
}

type UserProjectXref struct {
	User_id    string `json:"user_id"`
	Project_id string `json:"project_id"`
}

type Ticket struct {
	Id            string  `json:"id"`
	Priority      int     `json:"priority"`
	Title         string  `json:"title"`
	Story_points  int     `json:"story_points"`
	Description   *string `json:"description"`
	Status        string  `json:"status"`
	Created_by_id string  `json:"created_by_id"`
	Assignee_id   *string `json:"assignee_id"`
	Project_id    string  `json:"project_id"`
	Component_id  *string `json:"component_id"`
	Updated_by_id string  `json:"updated_by_id"`
	Created_at    string  `json:"created_at"`
	Updated_at    string  `json:"updated_at"`
}
