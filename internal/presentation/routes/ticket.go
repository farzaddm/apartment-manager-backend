package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

func SetUpTicketRoutes(publicG *gin.RouterGroup, managerG *gin.RouterGroup, handler *controller.TicketController) {
	publicG.POST("/tickets", handler.Create)
	publicG.GET("/tickets/:id", handler.GetByID)
	publicG.GET("/tickets/:id/fully", handler.GetFully)
	publicG.GET("/tickets", handler.List)
	publicG.PUT("/tickets/:id", handler.Update)
	publicG.GET("/tickets/me", handler.GetUserTickets)

	managerG.PATCH("/tickets/:id/status", handler.UpdateTicketStatus)
	managerG.DELETE("/tickets/:id", handler.Delete)
}
