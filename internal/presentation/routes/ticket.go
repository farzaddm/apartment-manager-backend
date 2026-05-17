package routes

import "github.com/gin-gonic/gin"

func SetUpTicketRoutes(g *gin.RouterGroup) {
	g.POST("/tickets")
	g.GET("/tickets/:id")
	g.GET("/tickets/:id/fully")
	g.GET("/tickets")
	g.PUT("tickets/:id")
	g.PATCH("tickets/:id/status")
	g.DELETE("tickets/:id")
}
