package db

type User struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Created_at string `json:"created_at"`
	Password   string `json:"password"`
}

type Project struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Leader_id   string `json:"leader_id"`
	Created_at  string `json:"created_at"`
}
