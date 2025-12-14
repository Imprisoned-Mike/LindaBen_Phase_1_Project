package models

import (
	"LindaBen_Phase_1_Project/internal/db"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Url string
}

// Save File details
func (file *File) Save() (*File, error) {
	if err := db.Db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

