package main

import (
	"fmt"
	"log"
	"os"

	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db.Connect(dsn) // your DB package would handle GORM connection

	// Setup Gin router
	r := gin.Default()

	// Example route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func seedData() {
	var roles = []models.Role{{RoleName: "admin", Description: "Administrator role"}, {RoleName: "customer", Description: "Authenticated customer role"}, {RoleName: "anonymous", Description: "Unauthenticated customer role"}}
	var user = []models.User{{Name: os.Getenv("ADMIN_NAME"), Email: os.Getenv("ADMIN_EMAIL"), Password: os.Getenv("ADMIN_PASSWORD"), RoleID: 1}}
	db.Db.Save(&roles)
	db.Db.Save(&user)
}

// run migration
func loadDatabase() {
	db.InitDb()
	db.Db.AutoMigrate(&models.Role{})
	db.Db.AutoMigrate(&models.User{})
	seedData()
}
