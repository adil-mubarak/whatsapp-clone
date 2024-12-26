package services

import (
	"net/http"
	"strconv"
	"whatsapp/db"
	"whatsapp/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateGroup(c *gin.Context) {
	var group models.Group

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group created successfully", "group": group})
}

func UpdateGroup(c *gin.Context) {
	groupID := c.Param("id")
	var group models.Group

	if err := db.DB.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group updated  successfully", "group": group})

}

func AddGroupMember(c *gin.Context) {
	groupIDParam := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var member models.GroupMember
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member.GroupID = groupID
	if err := db.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member to group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully added member to group", "group_member": member})
}

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
