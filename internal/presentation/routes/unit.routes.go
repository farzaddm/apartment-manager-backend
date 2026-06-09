package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

func SetUpUnitRoutes(managementG *gin.RouterGroup,
	handler *controller.UnitController) {
	managementG.POST("/apartments/:apartment_id/units", handler.Create)
	managementG.PUT("/apartments/units/:unit_id", handler.Update)
	managementG.PATCH("/apartments/units/:unit_id", handler.PopUser)
	managementG.DELETE("/apartments/units/:unit_id", handler.Delete)
	managementG.GET("/apartments/units/:unit_id", handler.GetByID)
	managementG.POST("/apartments/units/:unit_id/users", handler.PushUser)
}
