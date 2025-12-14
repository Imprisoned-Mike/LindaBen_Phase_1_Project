package handlers

import (
	Login "LindaBen_Phase_1_Project/internal/login"
	"LindaBen_Phase_1_Project/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Register user
func Register(context *gin.Context) {
	var input Login.Register

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.Users{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		RoleID:   3,
	}

	savedUser, err := user.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user": savedUser})

}

// User Login
func UserLogin(context *gin.Context) {
	var input Login.Login

	if err := context.ShouldBindJSON(&input); err != nil {
		var errorMessage string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationError := validationErrors[0]
			if validationError.Tag() == "required" {
				errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
			}
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	user, err := models.GetUserByName(input.Name)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidateUserPassword(input.Password)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"username": input.Name, "message": "Successfully logged in"})

}

// get all users
func GetUsers(c *gin.Context) {
	// Bind query parameters into UserFilterParams
	var filters models.UserFilterParams
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
	response, err := models.QueryUsers(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the paginated response
	c.JSON(http.StatusOK, response)
}

// get user by id
func GetUser(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var user models.Users
	err := models.GetUser(&user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, user)
}

// update user
func UpdateUser(c *gin.Context) {
	//var input models.Update
	var user models.Users
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.GetUser(&user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&user)
	err = models.UpdateUser(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}
