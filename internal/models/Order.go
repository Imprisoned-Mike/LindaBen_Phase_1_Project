package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Item     string
	Quantity int
	Unit     string  // "Kg", "L", "", etc.
	UnitCost float64 `json:"-"` //Should be hidden from frontend and just used in future for cost calculations

	DeliveryID uint
}
