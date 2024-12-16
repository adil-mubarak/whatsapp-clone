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
	

	router.Use(service.AuthMiddleWare())
	{
	router.POST("/profileUserName", services.ProfileUserName)
	router.POST("/profilePicture",services.ProfilePicture)
	}	
	return router
}
