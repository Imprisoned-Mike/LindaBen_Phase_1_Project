package models

import (
	"time"

	"gorm.io/gorm"
	"LindaBen_Phase_1_Project/internal/db"
)

// Delivery model
type Delivery struct {
	gorm.Model

	// One‐to‐many
	Orders []Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	BoxType      string    // e.g. type of box/packaging
	DeliveryDate time.Time // the scheduled or actual delivery datetime

	// Status
	Status string // e.g. "pending", "in_progress", "completed", "cancelled", etc.

	// Belongs-to: School
	SchoolID *uint
	School   *School `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET SET NULL;"`

	// Belongs-to: Vendor
	VendorID *uint
	Vendor   *Vendor `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Get all Deliveries
func GetAllDeliveries(Delivery *[]Delivery) (err error) {
	err = db.Db.Find(Delivery).Error
	if err != nil {
		return err
	}
	return nil
}

// Get Delivery by ID
func GetDeliveryByID(Delivery *Delivery, id uint) (err error) {
	err = db.Db.First(Delivery, id).Error
	if err != nil {
		return err
	}
	return nil
}

// Create Delivery
func CreateDelivery(Delivery *Delivery) (err error) {
	err = db.Db.Create(Delivery).Error
	if err != nil {
		return err
	}
	return nil
}

// Update Delivery
func UpdateDelivery(Delivery *Delivery) (err error) {
	err = db.Db.Save(Delivery).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete Delivery
func DeleteDelivery(Delivery *Delivery) (err error) {
	err = db.Db.Delete(Delivery).Error
	if err != nil {
		return err
	}
	return nil
}
