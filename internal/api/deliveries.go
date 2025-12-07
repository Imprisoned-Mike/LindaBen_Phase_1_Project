package api

import (
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// get all Deliveries
func GetDeliveries(context *gin.Context) {
	var deliveries []models.Delivery
	err := db.Db.Preload("Vendor").Preload("School").Preload("Orders").Find(&deliveries).Error
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, deliveries)
}

// get delivery by id
func GetDelivery(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var delivery models.Delivery
	err := db.Db.Preload("Vendor").Preload("School").Preload("Orders").First(&delivery, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
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
