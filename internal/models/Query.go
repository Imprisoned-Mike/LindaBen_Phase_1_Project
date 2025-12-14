package models

import (
	"LindaBen_Phase_1_Project/internal/db"

	"github.com/gin-gonic/gin"
)

func QueryUsers(
	page int,
	pageSize int,
	sortBy string,
	sortOrder string,
	expand []string,
	c *gin.Context,
) ([]Users, int, int, error) {

	var users []Users
	var total int64
	var totalUnfiltered int64

	query := db.Db.Model(&Users{}).
		Joins("UserRole")

	// ---- Filters ----
	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"users.name LIKE ? OR users.email LIKE ? OR users.phone LIKE ?",
			like, like, like,
		)
	}

	if email := c.Query("email"); email != "" {
		query = query.Where("users.email LIKE ?", "%"+email+"%")
	}

	if name := c.Query("name"); name != "" {
		query = query.Where("users.name LIKE ?", "%"+name+"%")
	}

	if role := c.Query("role"); role != "" {
		query = query.Where("roles.role_name = ?", role)
	}

	// Count BEFORE pagination
	query.Count(&total)

	// Count unfiltered (optional but useful)
	db.Db.Model(&Users{}).Count(&totalUnfiltered)

	// ---- Sorting ----
	order := sortBy
	if sortOrder == "desc" {
		order += " DESC"
	}
	query = query.Order(order)

	// ---- Pagination ----
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// ---- Expand relations ----
	for _, field := range expand {
		query = query.Preload(field)
	}

	err := query.Find(&users).Error
	return users, int(total), int(totalUnfiltered), err
}
