package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (u *UserController) Update(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "user_not_found_in_session", nil)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation_failed", err)
		return
	}

	err := u.userService.Update(c, req, userId.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res := dto.UpdateProfileResponse{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Gender:    req.Gender,
	}

	response.Success(c, http.StatusOK, "profile_updated_successfully", res)
}

func (u *UserController) ChangePassword(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "user_not_found_in_session", nil)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation_failed", err)
		return
	}

	err := u.userService.ChangePassword(c, req, userId.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "profile_updated_successfully", nil)
}

func (u *UserController) Delete(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "user_not_found_in_session", nil)
		return
	}

	err := u.userService.Delete(c.Request.Context(), userId.(string))
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

func (u *UserController) GetMe(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "user_not_found_in_session", nil)
		return
	}

	user, err := u.userService.GetById(c.Request.Context(), userId.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "profile_retrieved_successfully", user)
}

func (u *UserController) SetProfileImage(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user_not_found_in_session", nil)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "form_file_failed", err)
		return
	}

	extension := filepath.Ext(file.Filename)
	uniqueFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().UnixNano(), extension)

	uploadDir := "./uploads/profiles/"
	filePath := filepath.Join(uploadDir, uniqueFileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed_to_save_image_to_server_file_system", err)
		return
	}

	err = u.userService.SetProfileImage(c.Request.Context(), fmt.Sprintf("%v", userId), filePath)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed_to_save_image", err)
		return
	}

	image_addr := &dto.UploadImageResponse{ProfileImageURL: &filePath}

	response.Success(c, http.StatusOK, "profile_updated_successfully", image_addr)
}

func (u *UserController) CheckPhoneNumber(c *gin.Context) {
	var req dto.CheckUserPhoneNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation_failed", err)
		return
	}

	ex, err := u.userService.ExistSPhoneNumber(c, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "profile_updated_successfully", ex)
}
