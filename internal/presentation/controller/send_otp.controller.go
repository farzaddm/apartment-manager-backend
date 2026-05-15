package controller

import (
	"net/http"

	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	sendOTPUseCase *service.SendOtpService
}

func NewAuthHandler(sendOTP *service.SendOtpService) *AuthHandler {
	return &AuthHandler{
		sendOTPUseCase: sendOTP,
	}
}

type SendOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (h *AuthHandler) SendOTPHandler(c *gin.Context) {
	var req SendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	err := h.sendOTPUseCase.Execute(req.Phone)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to send otp", err)
		return
	}
	response.Success(c, http.StatusOK, "success", nil)
}
