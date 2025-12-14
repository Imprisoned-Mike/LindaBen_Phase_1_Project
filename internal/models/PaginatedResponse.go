package models

type PaginatedResponse[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Total           int `json:"total"`
	TotalUnfiltered int `json:"totalUnfiltered"`
	Page            int `json:"page"`
	PageSize        int `json:"pageSize"`
	TotalPages      int `json:"totalPages"`
}
