package services

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"whatsapp/db"
	"whatsapp/models"
	"whatsapp/pkg/service"

	"github.com/gin-gonic/gin"
)

type CreateStatusRequest struct {
	MediaURL string `json:"media_url"`
}

func CreateStatus(c *gin.Context) {
	claims, exist := c.Get("id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not Autherzide"})
		return
	}
	userClaims, ok := claims.(*service.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token data"})
		return
	}

	userID := userClaims.ID

	var req CreateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := models.StatusUpdate{
		UserID:    userID,
		MediaURL:  req.MediaURL,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := db.DB.Create(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create status"})
		return
	}

	c.JSON(http.StatusOK, status)
}

func ViewStatus(c *gin.Context) {
	calims, exist := c.Get("id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not Autherized"})
		return
	}

	userClaims, ok := calims.(*service.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token data"})
		return
	}

	userID := userClaims.ID

	var status []models.StatusUpdate
	if err := db.DB.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Preload("User").Find(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrive statuses"})
		return
	}

	c.JSON(http.StatusOK, status)
}

func UploadFileToStatus(c *gin.Context){
	claims,exist := c.Get("id")
	if !exist{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"user not authorized"})
		return
	}

	userClaims,ok := claims.(*service.Claims)
	if !ok{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"invalid token data"})
		return
	}

	userID := userClaims.ID

	file,err := c.FormFile("file")
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"file is required"})
		return
	}

	uploadDir := "./uploads/status"
	if _,err := os.Stat(uploadDir); os.IsNotExist(err){
		err := os.MkdirAll(uploadDir,os.ModePerm)
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"failed to create upload directionary"})
			return
		}
	}

	fileName := fmt.Sprintf("%d_%s",userID,file.Filename)
	filePath := filepath.Join(uploadDir,fileName)

	if err := 	c.SaveUploadedFile(file,filePath); err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"failed to save status"})
		return
	}

	status := models.StatusUpdate{
		UserID: userID,
		MediaURL: filePath,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := db.DB.Save(&status).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"failed to save status to database"})
		return
	}

	c.JSON(http.StatusOK,status)

}