package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	var err error
	Db, err = gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return nil
	}

	log.Println("Successfully connected to the database")
	return Db
}
