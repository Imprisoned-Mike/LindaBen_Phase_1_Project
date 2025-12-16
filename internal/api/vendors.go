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

// get all vendors
func GetVendors(c *gin.Context) {
	// Bind query parameters into VendorFilterParams
	var filters models.VendorFilterParams
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default pagination if not provided
	if filters.Page == nil {
		defaultPage := 1
		filters.Page = &defaultPage
	}
	if filters.PageSize == nil {
		defaultPageSize := 10
		filters.PageSize = &defaultPageSize
	}

	// Set default sorting if not provided
	if filters.SortBy == nil {
		defaultSort := "id"
		filters.SortBy = &defaultSort
	}
	if filters.SortOrder == nil {
		defaultOrder := "asc"
		filters.SortOrder = &defaultOrder
	}

	// Call QueryUsers with correct struct
	response, err := models.QueryVendors(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the paginated response
	c.JSON(http.StatusOK, response)
}

// get vendor by id
func GetVendor(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var vendor models.Vendor
	err := db.Db.Preload("Contact").First(&vendor, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, vendor)
}

// update vendor
func UpdateVendor(c *gin.Context) {
	//var input models.Update
	var vendor models.Vendor
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.GetVendorByID(&vendor, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&vendor)
	err = models.UpdateVendor(&vendor)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, vendor)
}

func DeleteVendor(c *gin.Context) {
	var vendor models.Vendor
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.GetVendorByID(&vendor, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = models.DeleteVendor(&vendor)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, vendor)
}

func CreateVendor(c *gin.Context) {
	var vendor models.Vendor

	if err := c.ShouldBindJSON(&vendor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.CreateVendor(&vendor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, vendor)
}
