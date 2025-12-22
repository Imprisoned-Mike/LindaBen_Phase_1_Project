package models

import (
	"time"

	"LindaBen_Phase_1_Project/internal/db"
)

// Delivery model
type Delivery struct {
	Model

	// One‐to‐many
	Orders []Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Contract string // hold, completed

	BoxType       string // standard, premium, standard-holiday, premium-holiday
	ScheduledFrom time.Time
	ScheduledTo   time.Time

	// Belongs-to: School
	SchoolID *uint
	School   *School `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Notes string
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
