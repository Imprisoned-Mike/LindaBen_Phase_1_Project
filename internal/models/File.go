package models

import (
	"LindaBen_Phase_1_Project/internal/db"

	"gorm.io/gorm"
)

type File struct {
	Model
	// Path is internal
	Path string `json:"-" gorm:"not null"`
	Url  string `json:"url" gorm:"-"`
}

// Save File details
func (file *File) Save() (*File, error) {
	if err := db.Db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

// Hook that Generate URL for file
func (file *File) AfterFind(tx *gorm.DB) (err error) {
	file.Url = "/uploads/" + file.Path
	return
}
