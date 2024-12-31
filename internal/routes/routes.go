package routs

import (
	"whatsapp/pkg/service"
	services "whatsapp/services"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/request-otp", services.RequestOTP)
	router.POST("/verify-otp", services.VerifyOTP)
	router.POST("/refresh-token", services.RefreshToken)

	router.Use(service.AuthMiddleWare())
	{
		router.POST("/profileUserName", services.ProfileUserName)
		router.POST("/profilePicture", services.ProfilePicture)

		router.POST("/messages", services.SendMessage)
		router.GET("/messages", services.GetMessage)

		router.POST("/createGroup", services.CreateGroup)
		router.PUT("/updateGroup/:id", services.UpdateGroup)
		router.POST("/addgroupmember/:id", services.AddGroupMember)
		router.POST("/removegroupmember/:id/:user_id", services.RemoveGroupMember)
		router.POST("/adminassign/:id", services.AssignAdmin)

		router.GET("/webmsg", service.WebSocketHandler)
		go service.BroadCastMessages()
	}

	return router

}
