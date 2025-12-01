package models

import "gorm.io/gorm"

type Vendor struct {
	gorm.Model
	Name    string
	Address string
	Lat     float64
	Lon     float64

	Type string // "produce", "shelf_stable", "packaging"
}
