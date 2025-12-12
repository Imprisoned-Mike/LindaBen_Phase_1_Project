package models

import "gorm.io/gorm"

type PaginatedResponse[T any] struct {
	gorm.Model
	Data []T `json:"data"`
	Meta struct {
		Total           int `json:"total"`
		TotalUnfiltered int `json:"totalUnfiltered"`
		Page            int `json:"page"`
		PageSize        int `json:"pageSize"`
		TotalPages      int `json:"totalPages"`
	} `json:"meta"`
}
