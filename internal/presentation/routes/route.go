package routes

import (
	"apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/presentation/controller"
	"apartment-manager-backend/internal/presentation/middleware"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Auth      *controller.AuthController
	Profile   *controller.ProfileController
	Apartment *controller.ApartmentController
	Ticket    *controller.TicketController
}

func SetUpRouter(handler *Controllers, jwtSvc jwt.TokenServiceInterface) *gin.Engine {
	r := gin.New()

	r.POST("/send-otp", handler.Auth.SendOTP)
	r.POST("/verify-otp", handler.Auth.VerifyOTP)
	r.POST("/register", handler.Auth.Register)
	r.POST("/login", handler.Auth.Login)
	r.POST("/refresh", handler.Auth.Refresh)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSvc))
	{
		protected.POST("/logout", handler.Auth.Logout)
		protected.GET("/me", handler.Profile.GetMe)
		SetUpApartmentRoutes(protected, handler.Apartment)
		SetUpTicketRoutes(protected, handler.Ticket)
	}

	return r
}
