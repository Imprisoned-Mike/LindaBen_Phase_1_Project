package models

import (
	"time"

	"LindaBen_Phase_1_Project/internal/db"

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

// Get Order by ID
func GetOrderByID(Order *Order, id uint) (err error) {
	err = db.Db.First(Order, id).Error
	if err != nil {
		return err
	}
	return nil
}

// Update Order
func UpdateOrder(Order *Order) (err error) {
	err = db.Db.Updates(Order).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete Order
func DeleteOrder(Order *Order) (err error) {
	err = db.Db.Delete(Order).Error
	if err != nil {
		return err
	}
	return nil
}
