package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

func SetUpApartmentRoutes(publicG *gin.RouterGroup, adminG *gin.RouterGroup, handler *controller.ApartmentController) {
	adminG.POST("/apartments", handler.Create)
	adminG.PUT("/apartments/:apartment_id", handler.Update)
	adminG.DELETE("/apartments/:apartment_id", handler.Delete)

	publicG.GET("/apartments/:apartment_id", handler.GetByID)
	publicG.GET("/apartments/:apartment_id/users", handler.GetByIDWithUsers)
	publicG.GET("/apartments/:apartment_id/rules", handler.GetByIDWithRules)
	publicG.GET("/apartments/:apartment_id/announcements", handler.GetByIDWithAnnouncements)
	publicG.GET("/apartments/:apartment_id/invite-codes", handler.GetByIDWithInviteCodes)
}
