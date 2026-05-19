package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InviteCodeController struct {
	inviteCodeService *service.InviteCodeService
}

func NewInviteCodeController(inviteCodeService *service.InviteCodeService) *InviteCodeController {
	return &InviteCodeController{inviteCodeService: inviteCodeService}
}

func (c *InviteCodeController) Create(ctx *gin.Context) {
	actorUserID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusBadRequest, "user_not_found_in_session", nil)
		return
	}

	var req dto.CreateInviteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	invite, err := c.inviteCodeService.Create(ctx.Request.Context(), req, fmt.Sprintf("%v", actorUserID))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	inviteCode := dto.InviteCodeResponse{
		ID:          invite.ID.String(),
		ApartmentID: invite.ApartmentID.String(),
		UnitID:      invite.UnitID.String(),
		Code:        invite.Code,
		ExpiresAt:   invite.ExpiresAt,
	}

	response.Success(ctx, http.StatusCreated, "invite_code_created_successfully", inviteCode)
}

func (c *InviteCodeController) Validate(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusBadRequest, "user_not_found_in_session", nil)
		return
	}

	var req dto.ValidateInviteCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "validation_failed", err)
		return
	}

	err := c.inviteCodeService.Validate(ctx.Request.Context(), req, fmt.Sprintf("%v", userID))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "successfully_joined_the_apartment_building_and_assigned_to_the_unit", nil)
}
