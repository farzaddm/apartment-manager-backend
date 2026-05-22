package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

func SetUpTicketRoutes(publicG *gin.RouterGroup, managerG *gin.RouterGroup, handler *controller.TicketController) {
	publicG.GET("/tickets/:id", handler.GetByID)        
	publicG.GET("/tickets/:id/fully", handler.GetFully) 
	publicG.GET("/tickets", handler.List)               
	publicG.GET("/tickets/me", handler.GetUserTickets)
	publicG.POST("/tickets", handler.Create)     
	publicG.PUT("/tickets/:id", handler.Update)   
	publicG.DELETE("/tickets/:id", handler.Delete) 

	managerG.PATCH("/tickets/:id/status", handler.UpdateTicketStatus)

	publicG.POST("/tickets/:id/comments", handler.CreateComment)

}
