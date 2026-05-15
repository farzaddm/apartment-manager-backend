package controller

import (
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogoutController struct {
	logoutService *service.LogoutService
}

func NewLogoutController(logoutService *service.LogoutService) *LogoutController {
	return &LogoutController{logoutService: logoutService}
}

func (h *LogoutController) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user_not_found_in_session", nil)
		return
	}

	err := h.logoutService.Execute(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "logged_out_successfully", nil)
}
