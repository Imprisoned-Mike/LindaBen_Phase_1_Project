package models

import (
	"LindaBen_Phase_1_Project/internal/db"
	"os"
	"strings"

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

func GetUrl(file *File) string {
	baseUrl := os.Getenv("BASE_URL")
	baseUrl = strings.TrimRight(baseUrl, "/")

	return baseUrl + "/api/uploads/" + file.Path
}

func (file *File) AfterFind(tx *gorm.DB) (err error) {
	file.Url = GetUrl(file)
	return
}

func (file *File) AfterCreate(tx *gorm.DB) (err error) {
	file.Url = GetUrl(file)
	return
}
