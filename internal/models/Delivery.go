package models

import (
	"time"

	"LindaBen_Phase_1_Project/internal/db"
)

// Delivery model
type Delivery struct {
	Model

	// One‐to‐many
	Orders []Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"orders"`

	Contract string `json:"contract"` // hold, completed

	BoxType       string `json:"packageType"` // standard, premium, standard-holiday, premium-holiday
	ScheduledFrom time.Time
	ScheduledTo   time.Time `json:"scheduledAt"`

	// Belongs-to: School
	SchoolID *uint   `json:"schoolId"`
	School   *School `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"school"`

	Notes string `json:"notes"`
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
