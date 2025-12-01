package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	RoleName    string       `gorm:"uniqueIndex"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
