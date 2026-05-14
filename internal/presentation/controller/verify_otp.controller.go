package controller

import (
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VerifyController struct {
	verifyOTPUseCase *service.VerifyOTPService
}

func NewVerifyController(verifyOTPUseCase *service.VerifyOTPService) *VerifyController {
	return &VerifyController{verifyOTPUseCase: verifyOTPUseCase}
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func (h *VerifyController) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	out, err := h.verifyOTPUseCase.Execute(
		c.Request.Context(),
		req.Phone,
		req.Code,
	)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "verification failed", err)
		return
	}

	response.Success(c, http.StatusOK, out.Message, gin.H{
		"user":          out.User,
		"access_token":  out.AccessToken,
		"refresh_token": out.RefreshToken,
	})
}
