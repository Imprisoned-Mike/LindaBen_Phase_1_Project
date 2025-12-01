package models

import "gorm.io/gorm"

type School struct {
	gorm.Model
	Name    string
	Address string
	Lat     float64
	Lon     float64

	ContactID *uint
	Contact   *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

