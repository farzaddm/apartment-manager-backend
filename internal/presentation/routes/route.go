package routes

import (
	"apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/presentation/controller"
	"apartment-manager-backend/internal/presentation/middleware"

	"github.com/gin-contrib/cors"
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
	Comment      *controller.CommentController
}

func SetUpRouter(handler *Controllers, jwtSvc jwt.TokenServiceInterface) *gin.Engine {
	r := gin.New()

	//TODO : Remove This is Release Mode
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	// CORS
	r.Use(cors.Default())

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
		protected.PUT("/user", handler.User.Update)
		protected.DELETE("/user", handler.User.Delete)
		protected.POST("/user/profile-image/", handler.User.SetProfileImage)
		protected.POST("/invite-code/validate", handler.InviteCode.Validate)
		protected.GET("/tags", handler.Tag.List)

		adminOnly := protected.Group("/")
		adminOnly.Use(middleware.RolesAuthorize("admin"))

		managementGroup := protected.Group("/")
		managementGroup.Use(middleware.RolesAuthorize("admin", "manager"))

		SetUpApartmentRoutes(protected, adminOnly, handler.Apartment)
		SetUpTicketRoutes(protected, managementGroup, handler.Ticket)
		SetUpCommentRoutes(protected, handler.Comment)

		managementGroup.POST("/invite-code", handler.InviteCode.Create)
		managementGroup.POST("/tags", handler.Tag.Create)
		managementGroup.DELETE("/tags/:id", handler.Tag.Delete)

		managementGroup.POST("/apartments/:apartment_id/announcements", handler.Announcement.Create)
		protected.GET("/apartments/:apartment_id/announcements/:id", handler.Announcement.Get)
		managementGroup.PUT("/apartments/:apartment_id/announcements/:id", handler.Announcement.Update)
		managementGroup.DELETE("/apartments/:apartment_id/announcements/:id", handler.Announcement.Delete)
	}

	return r
}
