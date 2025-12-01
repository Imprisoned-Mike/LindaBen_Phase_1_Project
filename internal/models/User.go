package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string `gorm:"unique"`
	Phone    string
	Role     string // e.g. "admin", "school_admin:34", "vendor_admin:12"

	AvatarID *uint
	Avatar   *File `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
