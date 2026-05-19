package routes

import (
	"apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/presentation/controller"
	"apartment-manager-backend/internal/presentation/middleware"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Auth         *controller.AuthController
	User         *controller.UserController
	Apartment    *controller.ApartmentController
	Ticket       *controller.TicketController
	InviteCode   *controller.InviteCodeController
	Tag          *controller.TagController
	Announcement *controller.AnnouncementController
}

func SetUpRouter(handler *Controllers, jwtSvc jwt.TokenServiceInterface) *gin.Engine {
	r := gin.New()

	r.Static("/uploads", "./uploads")

	r.POST("/send-otp", handler.Auth.SendOTP)
	r.POST("/verify-otp", handler.Auth.VerifyOTP)
	r.POST("/register", handler.Auth.Register)
	r.POST("/login", handler.Auth.Login)
	r.POST("/refresh", handler.Auth.Refresh)

	r.GET("/user/:user_id", handler.User.GetById)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSvc))
	{
		protected.POST("/logout", handler.Auth.Logout)
		SetUpApartmentRoutes(protected, handler.Apartment)
		SetUpTicketRoutes(protected, handler.Ticket)

		protected.PUT("/user", handler.User.Update)
		protected.DELETE("/user", handler.User.Delete)
		protected.POST("/user/profile-image/", handler.User.SetProfileImage)

		protected.POST("/invite-code", handler.InviteCode.Create)
		protected.POST("/invite-code/validate", handler.InviteCode.Validate)

		r.POST("/tags", handler.Tag.Create)
		r.GET("/tags", handler.Tag.List)
		r.DELETE("/tags/:id", handler.Tag.Delete)

		r.POST("/apartments/:apartment_id/announcements", handler.Announcement.Create)
		r.GET("/apartments/:apartment_id/announcements/:id", handler.Announcement.Get)
		r.PUT("/apartments/:apartment_id/announcements/:id", handler.Announcement.Update)
		r.DELETE("/apartments/:apartment_id/announcements/:id", handler.Announcement.Delete)
	}

	return r
}
