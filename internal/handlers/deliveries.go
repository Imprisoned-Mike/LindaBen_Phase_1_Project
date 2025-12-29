package handlers

import (
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterDeliveryRoutes registers delivery routes
func RegisterDeliveryRoutes(r *gin.RouterGroup) {
	r.GET("", GetDeliveries)
	r.GET("/:id", GetDelivery)
	r.POST("", CreateDelivery)
	r.PUT("/:id", UpdateDelivery)
	r.DELETE("/:id", DeleteDelivery)
	r.POST("/:delivery_id/orders", AddOrderToDelivery)
	r.DELETE("/:id/orders/:order_id", RemoveOrderFromDelivery)
}

// get all Deliveries
func GetDeliveries(context *gin.Context) {
	var filters models.DeliveryFilterParams
	if err := context.ShouldBindQuery(&filters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Always preload orders
	filters.Expand = append(filters.Expand, "order")

	response, err := models.QueryDeliveries(filters)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, response)
}

// get delivery by id
func GetDelivery(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))

	// Bind query params
	var expand []string
	if e := context.QueryArray("expand"); len(e) > 0 {
		expand = e
	}

	var delivery models.Delivery
	query := db.Db.Model(&models.Delivery{})

	// Preload associated data if requested
	for _, field := range expand {
		switch field {
		case "School":
			query = query.Preload("School")
		case "order":
			query = query.Preload("Orders.Vendor")
		}
	}

	// Get delivery by ID
	if err := query.First(&delivery, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, delivery)
}

// update delivery
func UpdateDelivery(c *gin.Context) {
	//var input models.Update
	var delivery models.Delivery
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.GetDeliveryByID(&delivery, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&delivery)
	err = models.UpdateDelivery(&delivery)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, delivery)
}

func DeleteDelivery(c *gin.Context) {
	var delivery models.Delivery
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.GetDeliveryByID(&delivery, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = models.DeleteDelivery(&delivery)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, delivery)
}

func CreateDelivery(c *gin.Context) {
	var delivery models.Delivery

	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.CreateDelivery(&delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, delivery)
}
