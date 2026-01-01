package handlers

import (
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/models"
	"LindaBen_Phase_1_Project/internal/util"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

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

		for _, delivery := range response.Data {
			filteredOrders := []models.Order{}
			for _, order := range delivery.Orders {
				// Check if order's vendor ID is in the list of vendor IDs
				if slices.Contains(vendorIDs, uint(*order.VendorID)) {
					filteredOrders = append(filteredOrders, order)
				}
			}
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
