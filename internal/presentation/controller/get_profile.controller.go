package controller

import (
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	profileService *service.ProfileService
}

func NewProfileController(profileService *service.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
	}
}

func (h *ProfileController) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user_not_found_in_session", nil)
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "profile_retrieved_successfully", profile)
}
