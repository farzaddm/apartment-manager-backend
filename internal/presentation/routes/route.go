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
	Comment      *controller.CommentController
	Rule         *controller.RuleController
	Poll         *controller.PollController
	Unit         *controller.UnitController
}

func SetUpRouter(handler *Controllers, jwtSvc jwt.TokenServiceInterface) *gin.Engine {
	r := gin.New()

	//TODO : Comment This is Release Mode
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS
	// r.Use(cors.Default())
	r.Use(corsMiddleware())

	r.Static("/uploads", "./uploads")

	r.POST("/send-otp", handler.Auth.SendOTP)
	r.POST("/verify-otp", handler.Auth.VerifyOTP)
	r.POST("/register", handler.Auth.Register)
	r.POST("/login", handler.Auth.Login)
	r.POST("/refresh", handler.Auth.Refresh)

	r.POST("/auth/check-phone", handler.User.CheckPhoneNumber)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSvc))
	{
		protected.POST("/logout", handler.Auth.Logout)
		protected.PUT("/user/me", handler.User.Update)
		protected.DELETE("/user/me", handler.User.Delete)
		protected.POST("/user/profile-image/", handler.User.SetProfileImage)
		protected.GET("/user/me", handler.User.GetMe)
		protected.PATCH("/user/me/password", handler.User.ChangePassword)
		protected.GET("/user/:user_id", handler.User.GetById)
		protected.POST("/invite-code/validate", handler.InviteCode.Validate)
		protected.GET("/tags", handler.Tag.List)

		adminOnly := protected.Group("/")
		adminOnly.Use(middleware.RolesAuthorize("admin"))

		managementGroup := protected.Group("/")
		managementGroup.Use(middleware.RolesAuthorize("admin", "manager"))

		SetUpApartmentRoutes(protected, managementGroup, adminOnly, handler.Apartment)
		SetUpTicketRoutes(protected, managementGroup, handler.Ticket)
		SetUpCommentRoutes(protected, handler.Comment)
		SetUpUnitRoutes(managementGroup, handler.Unit)

		managementGroup.POST("/invite-code", handler.InviteCode.Create)
		managementGroup.POST("/tags", handler.Tag.Create)
		managementGroup.DELETE("/tags/:id", handler.Tag.Delete)

		managementGroup.POST("/apartments/:apartment_id/announcements", handler.Announcement.Create)
		protected.GET("/apartments/:apartment_id/announcements/:id", handler.Announcement.Get)
		managementGroup.PUT("/apartments/:apartment_id/announcements/:id", handler.Announcement.Update)
		managementGroup.DELETE("/apartments/:apartment_id/announcements/:id", handler.Announcement.Delete)

		managementGroup.POST("/apartments/:apartment_id/rules", handler.Rule.Create)
		protected.GET("/apartments/:apartment_id/rules", handler.Rule.List)
		protected.GET("/apartments/:apartment_id/rules/:id", handler.Rule.Get)
		managementGroup.PUT("/apartments/:apartment_id/rules/:id", handler.Rule.Update)
		managementGroup.DELETE("/apartments/:apartment_id/rules/:id", handler.Rule.Delete)

		managementGroup.POST("/apartments/:apartment_id/polls", handler.Poll.Create)
		managementGroup.DELETE("/apartments/:apartment_id/polls/:poll_id", handler.Poll.Delete)
		protected.GET("/apartments/:apartment_id/polls", handler.Poll.List)
		protected.GET("/apartments/:apartment_id/polls/:poll_id", handler.Poll.GetDetails)
		protected.POST("/apartments/:apartment_id/polls/:poll_id/votes", handler.Poll.CastVote)
		protected.DELETE("/apartments/:apartment_id/polls/:poll_id/votes", handler.Poll.RevokeVote)
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE , PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
