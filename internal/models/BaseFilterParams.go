package models

type PaginationParams struct {
	Page     *int `form:"page"`
	PageSize *int `form:"pageSize"`
}

type SortParams struct {
	SortBy    *string `form:"sortBy"`
	SortOrder *string `form:"sortOrder"`
}

type BaseFilterParams struct {
	PaginationParams
	SortParams
	Expand []string `form:"expand"`
}
