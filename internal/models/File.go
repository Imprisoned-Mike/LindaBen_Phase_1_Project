package models

import (
	"LindaBen_Phase_1_Project/internal/db"
)

type File struct {
	Model
	Url string
}

// Save File details
func (file *File) Save() (*File, error) {
	if err := db.Db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

