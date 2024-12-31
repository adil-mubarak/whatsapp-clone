package services

import (
	"net/http"
	"whatsapp/db"
	"whatsapp/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RemoveGroupMember(c *gin.Context) {
	groupID := c.Param("id")
	userID := c.Param("user_id")
	if err := db.DB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMember{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove group member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed successfully"})
}

func AssignAdmin(c *gin.Context) {
	groupID := c.Param("id")
	var payload struct {
		UserID uint `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var groupMember models.GroupMember
	if err := db.DB.Where("group_id = ? AND user_id = ?", groupID, payload.UserID).First(&groupMember).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The user is not a member of the group"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify goup member"})
		return
	}

	if err := db.DB.Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, payload.UserID).Update("is_admin", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin assigned successfully"})
}

func RevokeAdmin(c *gin.Context) {
	groupID := c.Param("id")
	var payload struct {
		UserID uint `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var groupMember models.GroupMember
	if err := db.DB.Where("group_id = ? AND user_id = ?", groupID, payload.UserID).First(&groupMember).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The user is not a member of the group"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify the group member"})
		return
	}

	if err := db.DB.Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, payload.UserID).Update("is_admin", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin revoked successfully"})
}
