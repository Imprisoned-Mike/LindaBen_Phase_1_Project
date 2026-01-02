package main

import (
	"log"
	"os"

	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/handlers"
	"LindaBen_Phase_1_Project/internal/models"
	"LindaBen_Phase_1_Project/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, proceeding with defaults")
	}

	// Connect to SQLite DB
	database := db.InitDb("./mydatabase.db") // SQLite file is handled in db package
	if database == nil {
		log.Fatal("Failed to connect to database")
	}

	loadDatabase()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://rx.harvey-l.com") // Or use "*" for development
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Refresh-Token")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	auth := r.Group("/api/auth")
	handlers.RegisterAuthRoutes(auth)

	user := r.Group("/api/users", util.JWTAuth("admin"))
	handlers.RegisterUserRoutes(user)

	school := r.Group("/api/schools")
	handlers.RegisterSchoolRoutes(school)

	vendor := r.Group("/api/vendors")
	handlers.RegisterVendorRoutes(vendor)

	deliveries := r.Group("/api/deliveries")
	handlers.RegisterDeliveryRoutes(deliveries)

	order := r.Group("/api/orders")
	handlers.RegisterOrderRoutes(order)

	// Register logs routes
	deliveries.GET("/:id/logs", handlers.GetDeliveryLogs)
	order.GET("/:id/logs", handlers.GetOrderLogs)

	r.Static("/api/uploads", os.Getenv("UPLOAD_PATH"))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func seedData() {
	if os.Getenv("ADMIN_EMAIL") == "" {
		log.Println("Skipping seed data: ADMIN_EMAIL not set")
		return
	}

	var users = []models.User{{Name: os.Getenv("ADMIN_USERNAME"), Email: os.Getenv("ADMIN_EMAIL"), Password: os.Getenv("ADMIN_PASSWORD"), Roles: "admin"}}

	for _, u := range users {
		log.Printf("Seeding user: %s with email: %s and roles: %s", u.Name, u.Email, u.Roles)

		var existingUser models.User
		err := db.Db.Where("email = ?", u.Email).First(&existingUser).Error

		if err == nil {
			// User exists, update details
			// We set the plaintext password; the BeforeSave hook in the model will hash it.
			existingUser.Name = u.Name
			existingUser.Roles = u.Roles
			existingUser.Password = u.Password
			db.Db.Save(&existingUser)
		} else {
			// Create new user
			newUser := models.User{
				Name:     u.Name,
				Email:    u.Email,
				Password: u.Password,
				Roles:    u.Roles,
			}
			db.Db.Create(&newUser)
		}
	}
}

// run migration
func loadDatabase() {
	db.Db.AutoMigrate(&models.User{})
	db.Db.AutoMigrate(&models.School{})
	db.Db.AutoMigrate(&models.Vendor{})
	db.Db.AutoMigrate(&models.Delivery{})
	db.Db.AutoMigrate(&models.Order{})
	db.Db.AutoMigrate(&models.DeliveryChangeLog{})
	db.Db.AutoMigrate(&models.OrderChangeLog{})

	seedData()
}
