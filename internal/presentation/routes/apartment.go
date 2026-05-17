package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

func SetUpApartmentRoutes(g *gin.RouterGroup, handler *controller.ApartmentController) {
	g.POST("/apartments", handler.Create)
	g.GET("/apartments/:id", handler.GetByID)
	g.GET("/apartments/:id/users", handler.GetByIDWithUsers)
	g.GET("/apartments/:id/rules", handler.GetByIDWithRules)
	g.GET("/apartments/:id/announcements", handler.GetByIDWithAnnouncements)
	g.GET("/apartments/:id/invite-codes", handler.GetByIDWithInviteCodes)

	g.PUT("/apartments/:id", handler.Update)
	g.DELETE("/apartments/:id", handler.Delete)
}
