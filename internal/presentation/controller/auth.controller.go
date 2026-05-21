package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"apartment-manager-backend/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (h *AuthController) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	code, err := h.authService.SendOTP(req.Phone)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to send otp", err)
		return
	}

	response.Success(c, http.StatusOK, "success", gin.H{
		"otp": code,
	})
}

func (h *AuthController) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	out, err := h.authService.VerifyOTP(
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

func (h *AuthController) Register(c *gin.Context) {
	var req dto.RegisterInput

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

	out, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		var msg string
		switch err.Error() {
		case "email_already_exists":
			msg = "this email is already registered"
		case "username_already_exists":
			msg = "this username is taken"
		case "phone_not_verified_or_expired":
			msg = "phone verification failed or expired"
		default:
			msg = "registration failed"
		}

		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    msg,
			Errors:     []string{err.Error()},
		}
		res.SendResponse(c)
		return
	}

	response.Success(c, http.StatusCreated, "user registered successfully", out)
}

func (h *AuthController) Refresh(c *gin.Context) {
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

	out, err := h.authService.RefreshToken(c.Request.Context(), req)
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

	response.Success(c, http.StatusOK, "token_refreshed", out)
}

func (h *AuthController) Login(c *gin.Context) {
	var req dto.LoginInput

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

	out, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		var statusCode int
		var msg string

		switch err.Error() {
		case "invalid_credentials":
			statusCode = http.StatusUnauthorized
			msg = "invalid username or password"
		default:
			statusCode = http.StatusInternalServerError
			msg = "an unexpected error occurred"
		}

		res := &response.StandardResponse{
			Success:    false,
			StatusCode: statusCode,
			Message:    msg,
			Errors:     []string{err.Error()},
		}
		res.SendResponse(c)
		return
	}

	response.Success(c, http.StatusOK, "login_successful", out)
}

func (h *AuthController) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user_not_found_in_session", nil)
		return
	}

	err := h.authService.Logout(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "logged_out_successfully", nil)
}
