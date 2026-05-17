package routes

import "github.com/gin-gonic/gin"

func SetUpApartmentRoutes(g *gin.RouterGroup) {
	g.POST("/apartments")
	g.GET("/apartments/:id")
	g.GET("/apartments/:id/users")
	g.GET("/apartments/:id/rules")
	g.GET("/apartments/:id/announcements")
	g.GET("/apartments/:id/invite-codes")

	g.PUT("/apartments/:id")
	g.DELETE("/apartments/:id")
}
