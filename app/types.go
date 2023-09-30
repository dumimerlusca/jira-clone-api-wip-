package app

type PaginationMetadata struct {
	Page       int `json:"page"`
	TotalCount int `json:"totalCount"`
	Limit      int `json:"limit"`
}

type PaginatedResponse struct {
	Payload  any                `json:"payload"`
	Metadata PaginationMetadata `json:"metadata"`
}
