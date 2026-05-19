package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

func SetUpApartmentRoutes(g *gin.RouterGroup, handler *controller.ApartmentController) {
	g.POST("/apartments", handler.Create)
	g.GET("/apartments/:apartment_id", handler.GetByID)
	g.GET("/apartments/:apartment_id/users", handler.GetByIDWithUsers)
	g.GET("/apartments/:apartment_id/rules", handler.GetByIDWithRules)
	g.GET("/apartments/:apartment_id/announcements", handler.GetByIDWithAnnouncements)
	g.GET("/apartments/:apartment_id/invite-codes", handler.GetByIDWithInviteCodes)

	g.PUT("/apartments/:apartment_id", handler.Update)
	g.DELETE("/apartments/:apartment_id", handler.Delete)
}
