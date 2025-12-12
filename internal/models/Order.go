package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID       uint
	Item     string
	Quantity int
	Unit     string  // "Kg", "L", "", etc.
	UnitCost float64 `json:"-"` //Should be hidden from frontend and just used in future for cost calculations

	PackingTime  time.Time
	PurchaseTime time.Time

	DeliveryID uint

	VendorID *uint
	Vendor   Vendor

	isInternal bool // Use later in fetch api to differentiate between internal and external orders

	Status    string // "Pending", "Confirmed", "Completed", "Cancelled"
	Completed string

	Notes string
}
