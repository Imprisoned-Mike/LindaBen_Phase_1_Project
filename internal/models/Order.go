package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Item     string
	Quantity int
	Unit     string // "Kg", "L", "", etc.

	DeliveryID uint
}
