package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"apartment-manager-backend/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	loginService *service.LoginService
}

func NewLoginController(loginService *service.LoginService) *LoginController {
	return &LoginController{loginService: loginService}
}

func (h *LoginController) Login(c *gin.Context) {
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

	out, err := h.loginService.Execute(c.Request.Context(), req)
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
