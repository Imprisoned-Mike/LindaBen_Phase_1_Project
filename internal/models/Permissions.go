package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	PermName string `gorm:"uniqueIndex"`

	RoleID uint
}
