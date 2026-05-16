package controller

import (
	"apartment-manager-backend/internal/application/dto/user"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (u *UserController) Update(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		response.Error(c, http.StatusBadRequest, "user_id_required", nil)
		return
	}

	var req userdto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation_failed", err)
		return
	}

	err := u.userService.Update(c, req, userId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "profile_updated_successfully", nil)
}

func (u *UserController) Delete(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		response.Error(c, http.StatusBadRequest, "user_id_required", nil)
		return
	}

	err := u.userService.Delete(c.Request.Context(), userId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "user_deleted_successfully", nil)
}

func (u *UserController) GetById(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		response.Error(c, http.StatusUnauthorized, "user_id_not_found_in_params", nil)
		return
	}

	user, err := u.userService.GetById(c.Request.Context(), userId)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "profile_retrieved_successfully", user)
}
