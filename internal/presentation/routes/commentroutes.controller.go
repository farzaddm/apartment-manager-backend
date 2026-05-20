package routes

import (
	"apartment-manager-backend/internal/presentation/controller"

	"github.com/gin-gonic/gin"
)

func SetUpCommentRoutes(publicG *gin.RouterGroup, handler *controller.CommentController) {
	publicG.PATCH("/comments/:id", handler.UpdateBody)
	publicG.DELETE("/comments/:id", handler.Delete)
	publicG.GET("/comments/:id", handler.GetByID)
}
