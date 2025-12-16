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
func GetSchools(c *gin.Context) {
	// Bind query parameters into SchoolFilterParams
	var filters models.SchoolFilterParams
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
	response, err := models.QuerySchools(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the paginated response
	c.JSON(http.StatusOK, response)
}

// get school by id
func GetSchool(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))

	// Bind query params
	var expand []string
	if e := context.QueryArray("expand"); len(e) > 0 {
		expand = e
	}

	var school models.School
	query := db.Db.Model(&models.School{})

	// Preload avatar if requested
	for _, field := range expand {
		if field == "contact" {
			query = query.Preload("Contact")
		}
	}

	// Get school by ID
	if err := query.First(&school, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func DeleteSchool(c *gin.Context) {
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
	err = models.DeleteSchool(&school)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, school)
}

func CreateSchool(c *gin.Context) {
	var school models.School

	if err := c.ShouldBindJSON(&school); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.CreateSchool(&school); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, school)
}
