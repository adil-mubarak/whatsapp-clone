package services

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"whatsapp/db"
	"whatsapp/models"

	"github.com/gin-gonic/gin"
)

func ProfileUserName(c *gin.Context) {
	type Profile struct {
		UserName string `json:"user_name"`
		ProfilePicture string `json:"profile_picture"`
	}

	ID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unautherized"})
		return
	}

	var req Profile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := db.DB.Model(&models.User{}).Where("id = ?", ID).Updates(map[string]interface{}{
		"user_name":       req.UserName,
		"profile_picture": req.ProfilePicture,
	}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "user name updated successfully",
		"user_name": req.UserName,
		"profile_picture":req.ProfilePicture,
	})
}

func ProfilePicture(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	uploadDir := "./uploads/profile"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
			return
		}
	}

	fileName := fmt.Sprintf("%d_%s", user.ID, file.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	user.Profile_Picture = filePath
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile with image URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Profile picture added successfully",
		"user_name":       user.UserName,
		"profile_picture": user.Profile_Picture,
	})
}
