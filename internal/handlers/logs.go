package handlers

import (
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDeliveryLogs retrieves change logs for a specific delivery
func GetDeliveryLogs(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var logs []models.DeliveryChangeLog
	err := db.Db.Where("delivery_id = ?", id).Preload("ChangedByUser").Order("changed_at desc").Find(&logs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetOrderLogs retrieves change logs for a specific order
func GetOrderLogs(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Check permissions
	var logs []models.OrderChangeLog
	err := db.Db.Where("order_id = ?", id).Preload("ChangedByUser").Order("changed_at desc").Find(&logs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
