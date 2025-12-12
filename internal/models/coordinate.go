package models

import "gorm.io/gorm"

type Coordinate struct {
	gorm.Model
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	SchoolID  uint
}
