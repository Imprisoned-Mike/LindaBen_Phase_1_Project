package main

import (
	"log"
	"os"

	"LindaBen_Phase_1_Project/internal/api"
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/handlers"
	"LindaBen_Phase_1_Project/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
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

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://rx.harvey-l.com")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		c.Next()
	})

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://rx.harvey-l.com")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Status(200)
	})

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
		user.POST("", handlers.CreateUser)
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
		deliveries.DELETE("/:id/orders/:order_id", api.RemoveOrderFromDelivery)
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
	if os.Getenv("ADMIN_EMAIL") == "" {
		log.Println("Skipping seed data: ADMIN_EMAIL not set")
		return
	}

	var users = []models.Users{{Name: os.Getenv("ADMIN_USERNAME"), Email: os.Getenv("ADMIN_EMAIL"), Password: os.Getenv("ADMIN_PASSWORD"), Roles: "admin"},
		{Name: os.Getenv("SCHOOL_USERNAME"), Email: os.Getenv("SCHOOL_EMAIL"), Password: os.Getenv("SCHOOL_PASSWORD"), Roles: "school_admin"},
		{Name: os.Getenv("VENDOR_USERNAME"), Email: os.Getenv("VENDOR_EMAIL"), Password: os.Getenv("VENDOR_PASSWORD"), Roles: "vendor_admin"},
	}

	for _, u := range users {
		var existing models.Users
		result := db.Db.Where("email = ?", u.Email).First(&existing)

		if result.RowsAffected == 0 {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			newUser := models.Users{
				Name:     u.Name,
				Email:    u.Email,
				Password: string(hashedPassword),
				Roles:    u.Roles,
			}

			db.Db.Create(&newUser)
		}
	}

}

// run migration
func loadDatabase() {
	db.Db.AutoMigrate(&models.Users{})
	db.Db.AutoMigrate(&models.School{})
	db.Db.AutoMigrate(&models.Vendor{})
	db.Db.AutoMigrate(&models.Delivery{})
	db.Db.AutoMigrate(&models.Order{})

	seedData()
}
