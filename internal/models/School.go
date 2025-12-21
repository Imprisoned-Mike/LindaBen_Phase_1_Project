package models

import (
	"LindaBen_Phase_1_Project/internal/db"

	"gorm.io/gorm"
)

type School struct {
	gorm.Model // includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string
	Address    string
	Coordinate Coordinate `gorm:"embedded"`

	Contact *Users `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //Belongs to User
}

// Get all Schools
func GetAllSchools(School *[]School) (err error) {
	err = db.Db.Find(School).Error
	if err != nil {
		return err
	}
	return nil
}

// Get School by ID
func GetSchoolByID(School *School, id uint) (err error) {
	err = db.Db.First(School, id).Error
	if err != nil {
		return err
	}
	return nil
}

// Create School
func CreateSchool(School *School) (err error) {
	err = db.Db.Create(School).Error
	if err != nil {
		return err
	}
	return nil
}

// Update School
func UpdateSchool(School *School) (err error) {
	err = db.Db.Save(School).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete School
func DeleteSchool(School *School) (err error) {
	err = db.Db.Delete(School).Error
	if err != nil {
		return err
	}
	return nil
}
