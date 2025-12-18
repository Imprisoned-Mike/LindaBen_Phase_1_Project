package models

import (
	"LindaBen_Phase_1_Project/internal/db"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	RoleName    string `gorm:"uniqueIndex"`
	Description string `json:"description"`
}

// Create a role
func CreateRole(Role *Role) (err error) {
	err = db.Db.Create(Role).Error
	if err != nil {
		return err
	}
	return nil
}

// Get all roles
func GetRoles(Role *[]Role) (err error) {
	err = db.Db.Find(Role).Error
	if err != nil {
		return err
	}
	return nil
}

// Get role by id
func GetRole(Role *Role, id int) (err error) {
	err = db.Db.Where("id = ?", id).First(Role).Error
	if err != nil {
		return err
	}
	return nil
}

// Update role
func UpdateRole(Role *Role) (err error) {
	db.Db.Save(Role)
	return nil
}
