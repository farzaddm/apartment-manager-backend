package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	SendOTP *controller.AuthHandler
}

func SetUpRouter(handler *Controllers) *gin.Engine {

	r := gin.New()
	r.POST("/send-otp", handler.SendOTP.SendOTPHandler)

	return r
}
