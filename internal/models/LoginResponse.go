package models

import "gorm.io/gorm"

type LoginResponse struct {
	gorm.Model
	Token string `json:"token"`
	User  Users  `json:"user"`
}
