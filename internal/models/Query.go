package models

import (
	"LindaBen_Phase_1_Project/internal/db"
	"fmt"
	"math"
	"strings"
	"time"
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
	HasContact    *bool    `form:"hasContact"`
	ContactUserID *uint    `form:"contactUserId"`
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
	SchoolID      []uint   `form:"schoolId"`
	VendorID      []uint   `form:"vendorId"`
	ScheduledFrom *string  `form:"scheduledFrom"`
	ScheduledTo   *string  `form:"scheduledTo"`
	Contract      []string `form:"contract"`
	PackageType   []string `form:"packageType"`
	Status        []string `form:"status"`
	Page          *int     `form:"page"`
	PageSize      *int     `form:"pageSize"`
	SortBy        *string  `form:"sortBy"`
	SortOrder     *string  `form:"sortOrder"`
	Expand        []string `form:"expand"`
}

func QueryUsers(filters UserFilterParams) (PaginatedResponse[User], error) {
	var users []User
	query := db.Db.Model(&User{}).Preload("Avatar")

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

	if filters.Role != nil {
		query = query.Where("roles LIKE ?", "%"+*filters.Role+"%")
	}
	if filters.HasRole != nil {
		query = query.Where("roles LIKE ?", "%"+*filters.HasRole+"%")
	}
	// Note: EntityID filtering requires parsing the complex role string (e.g. "school_admin:123")
	// which is difficult to do reliably in a simple SQL WHERE clause without schema changes.

	// Total counts
	var total int64
	query.Count(&total)
	var totalUnfiltered int64
	db.Db.Model(&User{}).Count(&totalUnfiltered)

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
		return PaginatedResponse[User]{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := PaginatedResponse[User]{
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
		query = query.Where("LOWER(name) LIKE ? OR LOWER(address) LIKE ?", s, s)
	}
	if filters.HasContact != nil {
		if *filters.HasContact {
			query = query.Where("contact_id IS NOT NULL")
		} else {
			query = query.Where("contact_id IS NULL")
		}
	}
	if filters.ContactUserID != nil {
		query = query.Where("contact_id = ?", *filters.ContactUserID)
	}

	// TODO: filter by Role, EntityID, HasRole using RoleParsed

	// Total counts
	var total int64
	query.Count(&total)
	var totalUnfiltered int64
	db.Db.Model(&School{}).Count(&totalUnfiltered)

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
		query = query.Where("LOWER(name) LIKE ? OR LOWER(address) LIKE ?", s, s)
	}

	// TODO: filter by Role, EntityID, HasRole using RoleParsed

	// Total counts
	var total int64
	query.Count(&total)
	var totalUnfiltered int64
	db.Db.Model(&Vendor{}).Count(&totalUnfiltered)

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
	query := db.Db.Model(&Delivery{}).Preload("School")

	// Expand
	for _, field := range filters.Expand {
		if field == "orders" || strings.HasPrefix(field, "orders.") {
			query = query.Preload("Orders.Vendor")
		}
	}

	// Filters
	if filters.Search != nil && *filters.Search != "" {
		s := "%" + strings.ToLower(*filters.Search) + "%"
		query = query.Joins("LEFT JOIN schools ON schools.id = deliveries.school_id")
		query = query.Where("LOWER(deliveries.notes) LIKE ? OR LOWER(schools.name) LIKE ?", s, s)
	}
	if filters.ScheduledFrom != nil && *filters.ScheduledFrom != "" {
		from, err := time.Parse(time.RFC3339, *filters.ScheduledFrom)
		if err != nil {
			from, err = time.Parse("2006-01-02", *filters.ScheduledFrom)
			if err != nil {
				return PaginatedResponse[Delivery]{}, fmt.Errorf("invalid scheduledFrom: %w", err)
			}
		}
		query = query.Where("deliveries.scheduled_at >= ?", from)
	}
	if filters.ScheduledTo != nil && *filters.ScheduledTo != "" {
		to, err := time.Parse(time.RFC3339, *filters.ScheduledTo)
		if err != nil {
			to, err = time.Parse("2006-01-02", *filters.ScheduledTo)
			if err != nil {
				return PaginatedResponse[Delivery]{}, fmt.Errorf("invalid scheduledTo: %w", err)
			}
		}
		query = query.Where("deliveries.scheduled_at <= ?", to)
	}
	if len(filters.Contract) > 0 {
		query = query.Where("deliveries.contract IN ?", filters.Contract)
	}
	if len(filters.SchoolID) > 0 {
		query = query.Where("deliveries.school_id IN ?", filters.SchoolID)
	}
	if len(filters.PackageType) > 0 {
		query = query.Where("deliveries.package_type IN ?", filters.PackageType)
	}
	if len(filters.Status) > 0 {
		var parts []string
		for _, st := range filters.Status {
			s := strings.ToLower(strings.TrimSpace(st))
			if s == "" {
				continue
			}
			switch s {
			case "pending":
				parts = append(parts, "EXISTS (SELECT 1 FROM orders WHERE orders.delivery_id = deliveries.id AND LOWER(orders.status) = 'pending')")
			case "confirmed":
				parts = append(parts, "NOT EXISTS (SELECT 1 FROM orders WHERE orders.delivery_id = deliveries.id AND LOWER(orders.status) = 'pending') AND EXISTS (SELECT 1 FROM orders WHERE orders.delivery_id = deliveries.id AND LOWER(orders.status) = 'confirmed')")
			case "completed":
				parts = append(parts, "NOT EXISTS (SELECT 1 FROM orders WHERE orders.delivery_id = deliveries.id AND LOWER(orders.status) IN ('pending','confirmed')) AND EXISTS (SELECT 1 FROM orders WHERE orders.delivery_id = deliveries.id AND LOWER(orders.status) = 'completed')")
			case "cancelled":
				parts = append(parts, "NOT EXISTS (SELECT 1 FROM orders WHERE orders.delivery_id = deliveries.id AND LOWER(orders.status) IN ('pending','confirmed','completed')) AND EXISTS (SELECT 1 FROM orders WHERE orders.delivery_id = deliveries.id AND LOWER(orders.status) = 'cancelled')")
			}
		}
		if len(parts) > 0 {
			query = query.Where("(" + strings.Join(parts, " OR ") + ") AND deliveries.contract != 'hold'")
		}
	}
	if len(filters.VendorID) > 0 {
		var parts []string
		for _, vid := range filters.VendorID {
			parts = append(parts, fmt.Sprintf("EXISTS (SELECT 1 FROM orders WHERE orders.delivery_id = deliveries.id AND orders.vendor_id = %d)", vid))
		}
		query = query.Where("(" + strings.Join(parts, " OR ") + ")")
	}

	// Total counts
	var total int64
	query.Count(&total)
	var totalUnfiltered int64
	db.Db.Model(&Delivery{}).Count(&totalUnfiltered)

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
	sortBy := "deliveries.id"
	sortOrder := "asc"
	if filters.SortBy != nil {
		switch *filters.SortBy {
		case "scheduledAt":
			sortBy = "deliveries.scheduled_to"
		case "packageType":
			sortBy = "deliveries.box_type"
		case "notes":
			sortBy = "deliveries.notes"
		case "contract":
			sortBy = "deliveries.contract"
		case "schoolId":
			sortBy = "deliveries.school_id"
		}
	}
	if filters.SortOrder != nil {
		o := strings.ToLower(*filters.SortOrder)
		if o == "asc" || o == "desc" {
			sortOrder = o
		}
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
