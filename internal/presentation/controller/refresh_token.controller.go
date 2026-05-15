package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"apartment-manager-backend/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshController struct {
	refreshService *service.RefreshTokenService
}

func NewRefreshController(refreshService *service.RefreshTokenService) *RefreshController {
	return &RefreshController{refreshService: refreshService}
}

func (h *RefreshController) Refresh(c *gin.Context) {
	var req dto.RefreshInput

	if err := c.ShouldBindJSON(&req); err != nil {
		errList := validator.ParseValidationErrors(err)

		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "validation_failed",
			Errors:     errList,
		}
		res.SendResponse(c)
		return
	}

	out, err := h.refreshService.Execute(c.Request.Context(), req)
	if err != nil {
		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusUnauthorized,
			Message:    "refresh_failed",
			Errors:     []string{err.Error()},
		}
		res.SendResponse(c)
		return
	}

	// در صورت موفقیت
	response.Success(c, http.StatusOK, "token_refreshed", out)
}
