package handlers

import (
	"LindaBen_Phase_1_Project/internal/db"
	"LindaBen_Phase_1_Project/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	// "github.com/goccy/go-yaml/token"
	"gorm.io/gorm"
)

// RegisterUserRoutes registers user management routes
func RegisterUserRoutes(r *gin.RouterGroup) {
	r.GET("", GetUsers)
	r.GET("/:id", GetUser)
	r.POST("/:id/avatar", UploadUserAvatar)
	r.POST("", CreateUser)
	r.PUT("/:id", UpdateUser)
	r.DELETE("/:id", DeleteUser)
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Roles    string `json:"roles" binding:"required"`
}

// Create user
func CreateUser(context *gin.Context) {
	var input RegisterRequest

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.Users{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: input.Password,
		Roles:    input.Roles,
	}

	savedUser, err := user.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, savedUser)

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

	// Bind query params
	var expand []string
	if e := context.QueryArray("expand"); len(e) > 0 {
		expand = e
	}

	var user models.Users
	query := db.Db.Model(&models.Users{})

	// Preload avatar if requested
	for _, field := range expand {
		if field == "avatar" {
			query = query.Preload("Avatar")
		}
	}

	// Get user by ID
	if err := query.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func DeleteUser(c *gin.Context) {
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
	err = models.DeleteUser(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// upload user avatar
func UploadUserAvatar(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var user models.Users
	err := models.GetUser(&user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	file, err := CreateFileFromUpload(fileHeader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = models.UploadUserAvatar(&user, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": file.Url})
}

func CreateFileFromUpload(fileHeader *multipart.FileHeader) (*models.File, error) {
	// ensure uploads dir exists
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
	path := filepath.Join("uploads", filename)

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	file := &models.File{
		Path: path,
	}

	return file.Save()
}
