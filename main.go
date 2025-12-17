package main

import (
	"fmt"
	"log"
	"os"

	"LindaBen_Phase_1_Project/internal/api"
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/handlers"
	"LindaBen_Phase_1_Project/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with defaults")
	}

	// Connect to SQLite DB
	database := db.InitDb() // SQLite file is handled in db package
	if database == nil {
		log.Fatal("Failed to connect to database")
	}

	loadDatabase()

	// Setup Gin router
	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", handlers.UserLogin)
		auth.POST("/logout", handlers.UserLogout)
	}

	user := r.Group("/api/users")
	{
		user.GET("", handlers.GetUsers)
		user.GET("/:id", handlers.GetUser)
		user.POST("/:id/avatar", handlers.UploadUserAvatar)
		user.POST("", handlers.Register)
		user.PUT("/:id", handlers.UpdateUser)
		user.DELETE("/:id", handlers.DeleteUser)
	}

	school := r.Group("/api/schools")
	{
		school.GET("", api.GetSchools)
		school.GET("/:id", api.GetSchool)
		school.POST("", api.CreateSchool)
		school.PUT("/:id", api.UpdateSchool)
		school.DELETE("/:id", api.DeleteSchool)
	}

	vendor := r.Group("/api/vendors")
	{
		vendor.GET("", api.GetVendors)
		vendor.GET("/:id", api.GetVendor)
		vendor.POST("", api.CreateVendor)
		vendor.PUT("/:id", api.UpdateVendor)
		vendor.DELETE("/:id", api.DeleteVendor)
	}

	deliveries := r.Group("/api/deliveries")
	{
		deliveries.GET("", api.GetDeliveries)
		deliveries.GET("/:id", api.GetDelivery)
		deliveries.POST("", api.CreateDelivery)
		deliveries.PUT("/:id", api.UpdateDelivery)
		deliveries.DELETE("/:id", api.DeleteDelivery)
		deliveries.POST("/:delivery_id/orders", api.AddOrderToDelivery)
		deliveries.DELETE("/:delivery_id/orders/:order_id", api.RemoveOrderFromDelivery)
	}

	order := r.Group("/api/orders")
	{
		order.GET("/:id", api.GetOrderByID)
		order.PUT("/:id", api.UpdateOrder)
		// order.GET("/:id/logs", api.GetOrderLogs) // WIP
		order.DELETE("/:id", api.DeleteOrder)
		// order.GET("items/search", api.SearchOrderItems) // WIP
	}

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
	var roles = []models.Role{{RoleName: "admin", Description: "Administrator role"}, {RoleName: "customer", Description: "Authenticated customer role"}, {RoleName: "anonymous", Description: "Unauthenticated customer role"}}
	var user = []models.Users{{Name: os.Getenv("ADMIN_NAME"), Email: os.Getenv("ADMIN_EMAIL"), Password: os.Getenv("ADMIN_PASSWORD"), RoleID: 1}}
	db.Db.Save(&roles)
	db.Db.Save(&user)
}

// run migration
func loadDatabase() {
	db.InitDb()
	db.Db.AutoMigrate(&models.Role{})
	db.Db.AutoMigrate(&models.Users{})
	db.Db.AutoMigrate(&models.School{})
	db.Db.AutoMigrate(&models.Vendor{})
	db.Db.AutoMigrate(&models.Delivery{})
	db.Db.AutoMigrate(&models.Order{})

	seedData()
	SeedDemoUser()
}

func SeedDemoUser() {
	demoUser := models.Users{
		Name:     "Demo User",
		Email:    "admin@logistics.com",
		Password: "admin123", // ideally hash this
		RoleID:   3,
	}

	_, err := demoUser.Save()
	if err != nil {
		fmt.Println("Failed to create demo user:", err)
	} else {
		fmt.Println("Demo user created successfully!")
	}
}
