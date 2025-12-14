package api

import (
	"LindaBen_Phase_1_Project/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// save file details
func SaveFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "file is required",
		})
		return
	}

	uploadPath := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), fileHeader.Filename)

	if err := c.SaveUploadedFile(fileHeader, uploadPath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	file := &models.File{
		Url: "/" + uploadPath,
	}

	savedFile, err := file.Save()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, savedFile)
}
