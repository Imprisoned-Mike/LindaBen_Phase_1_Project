package models

import (
	"LindaBen_Phase_1_Project/internal/db"
)

type School struct {
	Model      // includes ID, CreatedAt, UpdatedAt
	Name       string
	Address    string
	Coordinate Coordinate `gorm:"embedded"`

	ContactID *uint
	Contact   *Users `gorm:"foreignKey:ContactID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //Belongs to User
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
