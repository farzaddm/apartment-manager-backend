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
		response.Error(ctx, http.StatusBadRequest, "invalid request body", err)
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

	response.Success(ctx, http.StatusCreated, "invite_code created successfully", inviteCode)
}
