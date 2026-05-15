package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	SendOTP   *controller.AuthHandler
	VerifyOTP *controller.VerifyController
	Register  *controller.RegisterController
	Login     *controller.LoginController
	Refresh   *controller.RefreshController
}

func SetUpRouter(handler *Controllers) *gin.Engine {

	r := gin.New()
	r.POST("/send-otp", handler.SendOTP.SendOTPHandler)
	r.POST("/verify-otp", handler.VerifyOTP.VerifyOTP)
	r.POST("/register", handler.Register.Register)
	r.POST("/login", handler.Login.Login)
	r.POST("/refresh", handler.Refresh.Refresh)

	return r
}
