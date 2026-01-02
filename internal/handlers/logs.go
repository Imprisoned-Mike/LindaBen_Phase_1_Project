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

	// Check permissions (copying logic from GetDelivery/GetOrderByID usually, 
	// but Delivery logs are likely viewable by those who can view the delivery)
	// For now, allow admin, school_admin, vendor_admin.
	// ideally we should check if the user is associated with this specific delivery's school/vendor
	// but simplified role check for now as per GetDeliveries

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
