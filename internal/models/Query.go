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

type SchoolFilterParams struct {
	Search        *string  `form:"search"`
	hasContact    *bool    `form:"hasContact"`
	contactUserID *uint    `form:"contactUserId"`
	Page          *int     `form:"page"`
	PageSize      *int     `form:"pageSize"`
	SortBy        *string  `form:"sortBy"`
	SortOrder     *string  `form:"sortOrder"`
	Expand        []string `form:"expand"`
}

type VendorFilterParams struct {
	Search    *string  `form:"search"`
	Types     []string `form:"types"`
	Page      *int     `form:"page"`
	PageSize  *int     `form:"pageSize"`
	SortBy    *string  `form:"sortBy"`
	SortOrder *string  `form:"sortOrder"`
	Expand    []string `form:"expand"`
}

type DeliveryFilterParams struct {
	Search        *string  `form:"search"`
	schoolID      *uint    `form:"schoolId"`
	scheduledFrom *string  `form:"scheduledFrom"`
	scheduledTo   *string  `form:"scheduledTo"`
	contract      *string  `form:"contract"`
	packageType   []string `form:"types"`
	Page          *int     `form:"page"`
	PageSize      *int     `form:"pageSize"`
	SortBy        *string  `form:"sortBy"`
	SortOrder     *string  `form:"sortOrder"`
	Expand        []string `form:"expand"`
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

func QuerySchools(filters SchoolFilterParams) (PaginatedResponse[School], error) {
	var school []School
	query := db.Db.Model(&School{}).Preload("Contact")

	// Filters
	if filters.Search != nil && *filters.Search != "" {
		s := "%" + strings.ToLower(*filters.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(phone) LIKE ?", s, s, s)
	}
	if filters.hasContact != nil {
		if *filters.hasContact {
			query = query.Where("contact_id IS NOT NULL")
		} else {
			query = query.Where("contact_id IS NULL")
		}
	}

	if filters.contactUserID != nil {
		query = query.Where("contact_user_id = ?", *filters.contactUserID)
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
	if err := query.Find(&school).Error; err != nil {
		return PaginatedResponse[School]{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := PaginatedResponse[School]{
		Data: school,
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

func QueryVendors(filters VendorFilterParams) (PaginatedResponse[Vendor], error) {
	var vendor []Vendor
	query := db.Db.Model(&Vendor{}).Preload("Contact")

	// Filters
	if filters.Search != nil && *filters.Search != "" {
		s := "%" + strings.ToLower(*filters.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(phone) LIKE ?", s, s, s)
	}

	if len(filters.Types) > 0 {
		query = query.Where("type IN ?", filters.Types)
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
	if err := query.Find(&vendor).Error; err != nil {
		return PaginatedResponse[Vendor]{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := PaginatedResponse[Vendor]{
		Data: vendor,
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

func QueryDeliveries(filters DeliveryFilterParams) (PaginatedResponse[Delivery], error) {
	var delivery []Delivery
	query := db.Db.Model(&Delivery{}).Preload("Contact")

	// Filters
	if filters.Search != nil && *filters.Search != "" {
		s := "%" + strings.ToLower(*filters.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(phone) LIKE ?", s, s, s)
	}
	if filters.scheduledFrom != nil {
		query = query.Where("scheduled_from >= ?", *filters.scheduledFrom)
	}
	if filters.scheduledTo != nil {
		query = query.Where("scheduled_to <= ?", *filters.scheduledTo)
	}
	if filters.contract != nil {
		query = query.Where("contract = ?", *filters.contract)
	}
	if filters.schoolID != nil {
		query = query.Where("school_id = ?", *filters.schoolID)
	}
	if len(filters.packageType) > 0 {
		query = query.Where("package_type IN ?", filters.packageType)
	}

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
	if err := query.Find(&delivery).Error; err != nil {
		return PaginatedResponse[Delivery]{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := PaginatedResponse[Delivery]{
		Data: delivery,
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
