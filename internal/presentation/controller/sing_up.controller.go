package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"apartment-manager-backend/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	registerUseCase *service.RegisterService
}

func NewRegisterController(registerUseCase *service.RegisterService) *RegisterController {
	return &RegisterController{registerUseCase: registerUseCase}
}

func (h *RegisterController) Register(c *gin.Context) {
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

	out, err := h.registerUseCase.Execute(c.Request.Context(), req)
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
