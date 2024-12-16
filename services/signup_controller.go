package services

import (
	"net/http"
	"strings"
	"time"
	"whatsapp/db"
	"whatsapp/models"
	"whatsapp/pkg/service"

	"github.com/gin-gonic/gin"
)

func RequestOTP(c *gin.Context) {
	type Request struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp, err := service.SendOTP(req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}
	var user models.User
	result := db.DB.Where("phone_number = ?", req.PhoneNumber).First(&user)

	if result.RowsAffected == 0 {
		user = models.User{
			PhoneNumber: req.PhoneNumber,
			OTP:         otp,
			OTPExpiry:   time.Now().Add(5 * time.Minute),
			CreatedAT:   time.Now(),
			UpdateAt:    time.Now(),
		}

		if err := db.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else {
		user.OTP = otp
		user.OTPExpiry = time.Now().Add(5 * time.Minute)
		user.UpdateAt = time.Now()
		if err := db.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func VerifyOTP(c *gin.Context) {
	type VerifyRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		OTP         string `json:"otp" binding:"required"`
	}

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)

	var user models.User
	if err := db.DB.Where("phone_number = ?", req.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.OTP != req.OTP || time.Now().After(user.OTPExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	db.DB.Model(&user).Updates(models.User{OTP: "", OTPExpiry: time.Time{}})

	token, err := service.GenerateJWT(user.ID, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := service.RefreshJWT(user.ID, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "OTP verified successfully",
		"token":        token,
		"refreshToken": refreshToken,
	})
}

func RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	refreshClaims, err := service.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newAccessToken, err := service.GenerateJWT(refreshClaims.ID, refreshClaims.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating new access token"})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"accessToken":newAccessToken,
	})
}
