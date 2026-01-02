package handlers

import (
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/models"
	"LindaBen_Phase_1_Project/internal/util"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterDeliveryRoutes registers delivery routes
func RegisterDeliveryRoutes(r *gin.RouterGroup) {
	r.GET("", util.JWTAuth("admin", "school_admin", "vendor_admin"), GetDeliveries)
	r.GET("/:id", GetDelivery)
	r.POST("", CreateDelivery)
	r.PUT("/:id", UpdateDelivery)
	r.DELETE("/:id", DeleteDelivery)
	r.POST("/:delivery_id/orders", AddOrderToDelivery)
	r.DELETE("/:id/orders/:order_id", RemoveOrderFromDelivery)
	r.POST("/:delivery_id/notify", util.JWTAuth("admin"), SendDeliveryNotifications)
	r.GET("/:id/notify/recipients", util.JWTAuth("admin"), GetDeliveryNotificationRecipients)
}

// get all Deliveries
func GetDeliveries(context *gin.Context) {
	var filters models.DeliveryFilterParams
	if err := context.ShouldBindQuery(&filters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize list params (expects repeated query keys, not comma-separated) and always preload orders
	normalize := func(vals []string) []string {
		var out []string
		for _, v := range vals {
			v = strings.TrimSpace(v)
			if v != "" {
				out = append(out, v)
			}
		}
		return out
	}
	filters.Expand = normalize(filters.Expand)
	filters.Contract = normalize(filters.Contract)
	filters.PackageType = normalize(filters.PackageType)
	filters.Status = normalize(filters.Status)
	filters.Expand = append(filters.Expand, "orders")

	// Check if is not admin, but school admin
	user := util.CurrentUser(context)
	userRoles := models.ParseRoles(user.Roles)

	adminRoleIdx := slices.IndexFunc(userRoles, func(r models.RoleParsed) bool {
		return r.Role == "admin"
	})
	vendorAdminRoleIdx := slices.IndexFunc(userRoles, func(r models.RoleParsed) bool {
		return r.Role == "vendor_admin"
	})
	schoolAdminRoleIdx := slices.IndexFunc(userRoles, func(r models.RoleParsed) bool {
		return r.Role == "school_admin"
	})

	if adminRoleIdx == -1 {
		for _, role := range userRoles {
			if role.Role == "school_admin" {
				// Filter by school ID
				filters.SchoolID = append(filters.SchoolID, *role.EntityID)
			}
			if role.Role == "vendor_admin" {
				// Filter by vendor ID
				filters.VendorID = append(filters.VendorID, *role.EntityID)
			}
		}
	}

	response, err := models.QueryDeliveries(filters)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter orders for vendor admins
	if adminRoleIdx == -1 && schoolAdminRoleIdx == -1 && vendorAdminRoleIdx != -1 {
		vendorIDs := []uint{}
		for _, role := range userRoles {
			if role.Role == "vendor_admin" {
				vendorIDs = append(vendorIDs, *role.EntityID)
			}
		}

		for i := range response.Data {
			delivery := &response.Data[i]

			filteredOrders := []models.Order{}
			for _, order := range delivery.Orders {
				// Check if order's vendor ID is in the list of vendor IDs
				if order.VendorID != nil && slices.Contains(vendorIDs, uint(*order.VendorID)) {
					filteredOrders = append(filteredOrders, order)
				} else {
					fmt.Println("Excluding order ID:", order.ID, "with VendorID:", order.VendorID)
				}
			}

			delivery.Orders = filteredOrders
		}
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
		case "school":
			query = query.Preload("School")
		case "orders":
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

func SendDeliveryNotifications(context *gin.Context) {
	context.JSON(http.StatusCreated, gin.H{"message": "Notifications sent (mock)."})
}

func GetDeliveryNotificationRecipients(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))

	var delivery models.Delivery
	query := db.Db.Preload("School").Preload("Orders.Vendor").Model(&models.Delivery{})
	// Get delivery by ID
	if err := query.First(&delivery, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users := []models.User{}

	// school admins
	if delivery.School != nil {
		schoolAdmins, _ := models.GetUsersByRole(fmt.Sprintf("school_admin:%d", delivery.School.ID))
		if schoolAdmins != nil {
			users = append(users, *schoolAdmins...)
		}
	}

	// all orders' vendor contacts
	for _, order := range delivery.Orders {
		if order.Vendor != nil {
			vendorAdmins, _ := models.GetUsersByRole(fmt.Sprintf("vendor_admin:%d", order.Vendor.ID))
			if vendorAdmins != nil {
				users = append(users, *vendorAdmins...)
			}
		}
	}

	context.JSON(http.StatusOK, users)
}

// update delivery
func UpdateDelivery(c *gin.Context) {
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

	// Clone for comparison
	oldDelivery := delivery

	if err := c.BindJSON(&delivery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Compare and Log
	currentUser := util.CurrentUser(c)
	var logs []models.DeliveryChangeLog
	now := time.Now()

	// Helper to add log
	addLog := func(field, old, new string) {
		logs = append(logs, models.DeliveryChangeLog{
			DeliveryID:     uint(delivery.ID),
			ChangeByUserID: uint(currentUser.ID),
			ChangedAt:      now,
			FieldName:      field,
			OldValue:       old,
			NewValue:       new,
		})
	}

	if oldDelivery.Contract != delivery.Contract {
		addLog("contract", oldDelivery.Contract, delivery.Contract)
	}

	if oldDelivery.PackageType != delivery.PackageType {
		addLog("packageType", oldDelivery.PackageType, delivery.PackageType)
	}

	if oldDelivery.Notes != delivery.Notes {
		addLog("notes", oldDelivery.Notes, delivery.Notes)
	}

	// ScheduledAt
	oldTime := ""
	newTime := ""
	if oldDelivery.ScheduledAt != nil {
		oldTime = oldDelivery.ScheduledAt.Format(time.RFC3339)
	}
	if delivery.ScheduledAt != nil {
		newTime = delivery.ScheduledAt.Format(time.RFC3339)
	}
	if oldTime != newTime {
		addLog("scheduledAt", oldTime, newTime)
	}

	// SchoolID
	oldSchool := ""
	newSchool := ""
	if oldDelivery.SchoolID != nil {
		oldSchool = fmt.Sprintf("%d", *oldDelivery.SchoolID)
	}
	if delivery.SchoolID != nil {
		newSchool = fmt.Sprintf("%d", *delivery.SchoolID)
	}
	if oldSchool != newSchool {
		addLog("schoolId", oldSchool, newSchool)
	}

	// Save Logs
	if len(logs) > 0 {
		db.Db.Create(&logs)
	}

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
