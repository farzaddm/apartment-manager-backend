package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

func SetUpTicketRoutes(g *gin.RouterGroup, handler *controller.TicketController) {
	g.POST("/tickets", handler.Create)
	g.GET("/tickets/:id", handler.GetByID)
	g.GET("/tickets/:id/fully", handler.GetFully)
	g.GET("/tickets", handler.List)
	g.PUT("/tickets/:id", handler.Update)
	g.PATCH("/tickets/:id/status", handler.UpdateTicketStatus)
	g.DELETE("/tickets/:id", handler.Delete)
	g.GET("/tickets/me", handler.GetUserTickets)
}
