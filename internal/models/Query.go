package models

import (
	"LindaBen_Phase_1_Project/internal/db"
	"fmt"
	"math"
	"strings"
)

type UserFilterParams struct {
	Search    *string  `form:"search"`
	Role      *string  `form:"role"`
	EntityID  *string  `form:"entityId"`
	HasRole   *string  `form:"hasRole"`
	ID        *uint    `form:"id"`
	Email     *string  `form:"email"`
	Name      *string  `form:"name"`
	Page      *int     `form:"page"`
	PageSize  *int     `form:"pageSize"`
	SortBy    *string  `form:"sortBy"`
	SortOrder *string  `form:"sortOrder"`
	Expand    []string `form:"expand"`
}

func QueryUsers(filters UserFilterParams) (PaginatedResponse[Users], error) {
	var users []Users
	query := db.Db.Model(&Users{}).Preload("UserRole")

	// Filters
	if filters.Search != nil && *filters.Search != "" {
		s := "%" + strings.ToLower(*filters.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(phone) LIKE ?", s, s, s)
	}
	if filters.ID != nil {
		query = query.Where("id = ?", *filters.ID)
	}
	if filters.Email != nil {
		query = query.Where("email = ?", *filters.Email)
	}
	if filters.Name != nil {
		query = query.Where("name = ?", *filters.Name)
	}

	// TODO: filter by Role, EntityID, HasRole using RoleParsed

	// Total counts
	var total int64
	query.Count(&total)
	var totalUnfiltered int64
	db.Db.Model(&Users{}).Count(&totalUnfiltered)

	// Pagination
	page := 1
	pageSize := 10
	if filters.Page != nil {
		page = *filters.Page
	}
	if filters.PageSize != nil {
		pageSize = *filters.PageSize
	}
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Sorting
	sortBy := "id"
	sortOrder := "asc"
	if filters.SortBy != nil {
		sortBy = *filters.SortBy
	}
	if filters.SortOrder != nil {
		sortOrder = *filters.SortOrder
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Execute
	if err := query.Find(&users).Error; err != nil {
		return PaginatedResponse[Users]{}, err
	}

	// Add RoleParsed
	for i := range users {
		users[i].RoleParsed = ParseRole(users[i])
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := PaginatedResponse[Users]{
		Data: users,
		Meta: PaginationMeta{
			Total:           int(total),
			TotalUnfiltered: int(totalUnfiltered),
			Page:            page,
			PageSize:        pageSize,
			TotalPages:      totalPages,
		},
	}

	return response, nil

}
