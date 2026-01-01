package handlers

import (
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/models"
	"LindaBen_Phase_1_Project/internal/util"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterOrderRoutes registers order routes
func RegisterOrderRoutes(r *gin.RouterGroup) {
	r.GET("/:id", GetOrderByID)
	r.PUT("/:id", UpdateOrder, util.JWTAuth("admin"))
	r.DELETE("/:id", DeleteOrder, util.JWTAuth("admin"))
}

// get order by id
func GetOrderByID(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var order models.Order
	err := db.Db.Preload("Delivery").Preload("Vendor").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	allowedRoles := []string{"admin", fmt.Sprintf("school_admin:%d", order.Delivery.SchoolID)}

	if order.VendorID != nil {
		allowedRoles = append(allowedRoles, fmt.Sprintf("vendor_admin:%d", *order.VendorID))
	}

	if err := util.ValidateRoleJWT(context, allowedRoles...); err != nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}

	context.JSON(http.StatusOK, order)
}

// update order
func UpdateOrder(c *gin.Context) {
	//var input models.Update
	var order models.Order
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.GetOrderByID(&order, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&order)
	err = models.UpdateOrder(&order)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, order)
}

func DeleteOrder(c *gin.Context) {
	var order models.Order
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.GetOrderByID(&order, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = models.DeleteOrder(&order)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, order)
}

// Add Order to Delivery
func AddOrderToDelivery(c *gin.Context) {
	var delivery models.Delivery
	var order models.Order
	deliveryID, _ := strconv.Atoi(c.Param("delivery_id"))
	orderID, _ := strconv.Atoi(c.Param("order_id"))

	err := db.Db.First(&delivery, deliveryID).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		return
	}

	err = db.Db.First(&order, orderID).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	err = models.AddOrderToDelivery(&delivery, &order)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order added to delivery successfully"})
}

// Remove Order from Delivery
func RemoveOrderFromDelivery(c *gin.Context) {
	var delivery models.Delivery
	var order models.Order
	deliveryID, _ := strconv.Atoi(c.Param("delivery_id"))
	orderID, _ := strconv.Atoi(c.Param("order_id"))

	err := db.Db.First(&delivery, deliveryID).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		return
	}

	err = db.Db.First(&order, orderID).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	err = models.RemoveOrderFromDelivery(&delivery, &order)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order removed from delivery successfully"})
}
