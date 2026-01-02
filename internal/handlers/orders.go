package handlers

import (
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/models"
	"LindaBen_Phase_1_Project/internal/util"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterOrderRoutes registers order routes
func RegisterOrderRoutes(r *gin.RouterGroup) {
	r.GET("/:id", GetOrderByID)
	r.PUT("/:id", UpdateOrder, util.JWTAuth("admin"))
	r.DELETE("/:id", DeleteOrder, util.JWTAuth("admin"))
	r.GET("/:id/notify/recipients", util.JWTAuth("admin"), GetOrderNotificationRecipients)
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

func GetOrderNotificationRecipients(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var order models.Order

	err := db.Db.Preload("Vendor").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Send to vendor admins if vendor exists
	users := []models.User{}
	if order.Vendor != nil {
		vendorAdmins, _ := models.GetUsersByRole(fmt.Sprintf("vendor_admin:%d", order.Vendor.ID))
		if vendorAdmins != nil {
			users = append(users, *vendorAdmins...)
		}
	}

	context.JSON(http.StatusOK, users)
}

// update order
func UpdateOrder(c *gin.Context) {
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

	// Clone for comparison
	oldOrder := order

	if err := c.BindJSON(&order); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Compare and Log
	currentUser := util.CurrentUser(c)
	var logs []models.OrderChangeLog
	now := time.Now()

	addLog := func(field, old, new string) {
		logs = append(logs, models.OrderChangeLog{
			OrderID:        uint(order.ID),
			ChangeByUserID: uint(currentUser.ID),
			ChangedAt:      now,
			FieldName:      field,
			OldValue:       old,
			NewValue:       new,
		})
	}

	if oldOrder.Status != order.Status {
		addLog("status", oldOrder.Status, order.Status)
	}
	if oldOrder.Quantity != order.Quantity {
		addLog("quantity", strconv.Itoa(oldOrder.Quantity), strconv.Itoa(order.Quantity))
	}
	if oldOrder.Item != order.Item {
		addLog("item", oldOrder.Item, order.Item)
	}
	if oldOrder.UnitCost != order.UnitCost {
		addLog("unitPrice", fmt.Sprintf("%f", oldOrder.UnitCost), fmt.Sprintf("%f", order.UnitCost))
	}
	if oldOrder.Notes != order.Notes {
		addLog("notes", oldOrder.Notes, order.Notes)
	}
	if oldOrder.IsInternal != order.IsInternal {
		addLog("isInternal", strconv.FormatBool(oldOrder.IsInternal), strconv.FormatBool(order.IsInternal))
	}

	// VendorID
	oldVendor := ""
	newVendor := ""
	if oldOrder.VendorID != nil {
		oldVendor = fmt.Sprintf("%d", *oldOrder.VendorID)
	}
	if order.VendorID != nil {
		newVendor = fmt.Sprintf("%d", *order.VendorID)
	}
	if oldVendor != newVendor {
		addLog("vendorId", oldVendor, newVendor)
	}

	// Save logs
	if len(logs) > 0 {
		db.Db.Create(&logs)
	}

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
