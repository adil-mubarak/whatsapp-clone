package services

import (
	"net/http"
	"strconv"
	"whatsapp/db"
	"whatsapp/models"

	"github.com/gin-gonic/gin"
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

func GetGroups(c *gin.Context) {
	var groups []models.Group
	if err := db.DB.Preload("Members").Find(&groups); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func GetGroup(c *gin.Context) {
	groupID := c.Param("id")
	var group models.Group

	if err := db.DB.Preload("Members").First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group not found"})
		return
	}
	c.JSON(http.StatusOK, group)
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

func ListOfGroupMember(c *gin.Context) {
	groupID := c.Param("id")
	var members []models.GroupMember
	if err := db.DB.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}
