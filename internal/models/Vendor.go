package models

import (
	"LindaBen_Phase_1_Project/internal/db"

	"gorm.io/gorm"
)

type Vendor struct {
	gorm.Model // includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string
	Address    string
	Location   float64 // e.g. GPS coordinates

	Type string // "produce", "shelf_stable", "packaging"
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
