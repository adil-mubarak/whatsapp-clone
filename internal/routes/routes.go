package routs

import (
	"net/http"
	"whatsapp/pkg/service"
	services "whatsapp/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

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
		router.GET("/groups", services.GetGroups)
		router.GET("/group/:id", services.GetGroup)
		router.GET("/groupMemberes/:id", services.ListOfGroupMember)

		router.POST("/status", services.CreateStatus)
		router.GET("/status", services.ViewStatus)
		router.POST("/uploadstatus", services.UploadFileToStatus)

		http.HandleFunc("/offer", service.HandleOffer)
		// router.GET("/webmsg", service.WebSocketHandler)
		// go service.BroadCastMessages()
	}

	return router

}
