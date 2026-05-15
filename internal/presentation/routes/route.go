package routes

import (
	"apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/presentation/controller"
	"apartment-manager-backend/internal/presentation/middleware"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	SendOTP   *controller.AuthHandler
	VerifyOTP *controller.VerifyController
	Register  *controller.RegisterController
	Login     *controller.LoginController
	Refresh   *controller.RefreshController
	Logout    *controller.LogoutController
}

func SetUpRouter(handler *Controllers, jwtSvc jwt.TokenServiceInterface) *gin.Engine {

	r := gin.New()
	r.POST("/send-otp", handler.SendOTP.SendOTPHandler)
	r.POST("/verify-otp", handler.VerifyOTP.VerifyOTP)
	r.POST("/register", handler.Register.Register)
	r.POST("/login", handler.Login.Login)
	r.POST("/refresh", handler.Refresh.Refresh)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSvc))
	{
		protected.POST("/logout", handler.Logout.Logout)

	}

	return r
}
