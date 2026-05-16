package routes

import (
	"apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/presentation/controller"
	"apartment-manager-backend/internal/presentation/middleware"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Auth *controller.AuthController
	User *controller.UserController
}

func SetUpRouter(handler *Controllers, jwtSvc jwt.TokenServiceInterface) *gin.Engine {
	r := gin.New()

	r.POST("/send-otp", handler.Auth.SendOTP)
	r.POST("/verify-otp", handler.Auth.VerifyOTP)
	r.POST("/register", handler.Auth.Register)
	r.POST("/login", handler.Auth.Login)
	r.POST("/refresh", handler.Auth.Refresh)

	// dont need id if user is authenticated
	r.PUT("/user/:user_id", handler.User.Update)
	r.DELETE("/user/:user_id", handler.User.Delete)
	r.GET("/user/:user_id", handler.User.GetById)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSvc))
	{
		protected.POST("/logout", handler.Auth.Logout)
	}

	return r
}
