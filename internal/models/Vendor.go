package models

import (
	"LindaBen_Phase_1_Project/internal/db"
)

type Vendor struct {
	Model      // includes ID, CreatedAt, UpdatedAt
	Name       string `gorm:"unique"`
	Address    string
	Coordinate Coordinate `gorm:"embedded"`

	Type string // "produce", "shelf_stable", "packaging"

	ContactID *uint
	Contact   *Users `gorm:"foreignKey:ContactID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //Belongs to User
}

// Get all Vendors
func GetAllVendors(Vendor *[]Vendor) (err error) {
	err = db.Db.Find(Vendor).Error
	if err != nil {
		return err
	}
	return nil
}

// Get Vendor by ID
func GetVendorByID(Vendor *Vendor, id uint) (err error) {
	err = db.Db.First(Vendor, id).Error
	if err != nil {
		return err
	}
	return nil
}

// Create Vendor
func CreateVendor(Vendor *Vendor) (err error) {
	err = db.Db.Create(Vendor).Error
	if err != nil {
		return err
	}
	return nil
}

// Update Vendor
func UpdateVendor(Vendor *Vendor) (err error) {
	err = db.Db.Save(Vendor).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete Vendor
func DeleteVendor(Vendor *Vendor) (err error) {
	err = db.Db.Delete(Vendor).Error
	if err != nil {
		return err
	}
	return nil
}
