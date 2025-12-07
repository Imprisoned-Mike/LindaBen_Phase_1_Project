package api

import (
	"LindaBen_Phase_1_Project/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// get all vendors
func GetVendors(context *gin.Context) {
	var vendor []models.Vendor
	err := models.GetAllVendors(&vendor)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, vendor)
}

// get vendor by id
func GetVendor(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var vendor models.Vendor
	err := models.GetVendorByID(&vendor, uint(id))
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
