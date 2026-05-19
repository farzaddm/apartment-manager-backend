package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

func SetUpTicketRoutes(publicG *gin.RouterGroup, managerG *gin.RouterGroup, handler *controller.TicketController) {
	publicG.GET("/tickets/:id", handler.GetByID)        //TODO:CHECK Public/Private field + and is it managerG+
	publicG.GET("/tickets/:id/fully", handler.GetFully) //TODO:CHECK Public/Private field + and is it managerG+
	publicG.GET("/tickets", handler.List)               //TODO:Filter By Public/Private field situation + and is it managerG+
	publicG.GET("/tickets/me", handler.GetUserTickets)
	publicG.POST("/tickets", handler.Create)       //TODO: Check base user id with target user id +
	publicG.PUT("/tickets/:id", handler.Update)    //TODO: Check base user id with target user id +
	publicG.DELETE("/tickets/:id", handler.Delete) //TODO: Check base user id with target user id + and is it managerG+

	managerG.PATCH("/tickets/:id/status", handler.UpdateTicketStatus)
}
