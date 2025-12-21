package models

import (
	"LindaBen_Phase_1_Project/internal/db"
	"strconv"
)

type RoleParsed struct {
	Role     string //admin, school_admin, vendor_admin, user
	EntityID *string
}

func ParseRole(user Users) RoleParsed {
	var entityId *string

	if user.Roles == "school_admin" {
		var school School

		err := db.Db.Where("contact_id = ?", user.ID).First(&school).Error
		if err == nil {
			idStr := strconv.Itoa(int(school.ID))
			entityId = &idStr
		}
	}

	return RoleParsed{
		Role:     user.Roles,
		EntityID: entityId,
	}
}

//repeat for admin, vendor_admin, user
