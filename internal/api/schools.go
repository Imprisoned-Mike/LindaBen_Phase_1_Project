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

// get all schools
func GetSchools(context *gin.Context) {
	var school []models.School
	err := db.Db.Preload("Contact").Find(&school).Error
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, school)
}

// get school by id
func GetSchool(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var school models.School
	err := db.Db.Preload("Contact").First(&school, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, school)
}

// update school
func UpdateSchool(c *gin.Context) {
	//var input models.Update
	var school models.School
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.GetSchoolByID(&school, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&school)
	err = models.UpdateSchool(&school)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, school)
}
