package models

import "gorm.io/gorm"

type PaginationParams struct {
	gorm.Model
	Page     *int `form:"page"`
	PageSize *int `form:"pageSize"`
}

type SortParams struct {
	gorm.Model
	SortBy    *string `form:"sortBy"`
	SortOrder *string `form:"sortOrder"`
}

type BaseFilterParams struct {
	gorm.Model
	PaginationParams
	SortParams
	Expand []string `form:"expand"`
}
