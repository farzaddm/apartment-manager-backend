package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PollController struct {
	pollService service.PollService
}

func NewPollController(pollService service.PollService) *PollController {
	return &PollController{pollService: pollService}
}

func (c *PollController) Create(ctx *gin.Context) {
	apartmentIDStr := ctx.Param("apartment_id")
	apartmentID, err := uuid.Parse(apartmentIDStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	var req dto.CreatePollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	resp, err := c.pollService.CreatePoll(apartmentID, req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "poll_created_successfully", resp)
}

func (c *PollController) List(ctx *gin.Context) {
	apartmentIDStr := ctx.Param("apartment_id")
	apartmentID, err := uuid.Parse(apartmentIDStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	userIDVal, _ := ctx.Get("user_id")
	userIDStr, _ := userIDVal.(string)
	userID, _ := uuid.Parse(userIDStr)

	userApartmentIDVal, _ := ctx.Get("apartment_id")
	var userApartmentID *uuid.UUID = nil
	if userApartmentIDVal != nil {
		uaIDStr, _ := userApartmentIDVal.(string)
		if parsed, errParse := uuid.Parse(uaIDStr); errParse == nil {
			userApartmentID = &parsed
		}
	}

	userRoleVal, _ := ctx.Get("role")
	userRole, _ := userRoleVal.(string)

	resp, err := c.pollService.ListPolls(apartmentID, userID, userApartmentID, userRole)
	if err != nil {
		response.Error(ctx, http.StatusForbidden, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "polls_fetched_successfully", resp)
}

func (c *PollController) GetDetails(ctx *gin.Context) {
	pollIDStr := ctx.Param("poll_id")
	pollID, err := uuid.Parse(pollIDStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_poll_id", err)
		return
	}

	userIDVal, _ := ctx.Get("user_id")
	userIDStr, _ := userIDVal.(string)
	userID, _ := uuid.Parse(userIDStr)

	userApartmentIDVal, _ := ctx.Get("apartment_id")
	var userApartmentID *uuid.UUID = nil
	if userApartmentIDVal != nil {
		uaIDStr, _ := userApartmentIDVal.(string)
		if parsed, errParse := uuid.Parse(uaIDStr); errParse == nil {
			userApartmentID = &parsed
		}
	}

	userRoleVal, _ := ctx.Get("role")
	userRole, _ := userRoleVal.(string)

	resp, err := c.pollService.GetPollDetails(pollID, userID, userApartmentID, userRole)
	if err != nil {
		response.Error(ctx, http.StatusForbidden, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "poll_details_fetched_successfully", resp)
}

func (c *PollController) CastVote(ctx *gin.Context) {
	pollIDStr := ctx.Param("poll_id")
	pollID, err := uuid.Parse(pollIDStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_poll_id", err)
		return
	}

	userIDVal, _ := ctx.Get("user_id")
	userIDStr, _ := userIDVal.(string)
	userID, _ := uuid.Parse(userIDStr)

	userApartmentIDVal, _ := ctx.Get("apartment_id")
	var userApartmentID *uuid.UUID = nil
	if userApartmentIDVal != nil {
		uaIDStr, _ := userApartmentIDVal.(string)
		if parsed, errParse := uuid.Parse(uaIDStr); errParse == nil {
			userApartmentID = &parsed
		}
	}

	userRoleVal, _ := ctx.Get("role")
	userRole, _ := userRoleVal.(string)

	var req dto.VoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	err = c.pollService.CastVote(pollID, userID, userApartmentID, userRole, req.OptionID)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "vote_registered_successfully", nil)
}

func (c *PollController) RevokeVote(ctx *gin.Context) {
	pollIDStr := ctx.Param("poll_id")
	pollID, err := uuid.Parse(pollIDStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_poll_id", err)
		return
	}

	userIDVal, _ := ctx.Get("user_id")
	userIDStr, _ := userIDVal.(string)
	userID, _ := uuid.Parse(userIDStr)

	userApartmentIDVal, _ := ctx.Get("apartment_id")
	var userApartmentID *uuid.UUID = nil
	if userApartmentIDVal != nil {
		uaIDStr, _ := userApartmentIDVal.(string)
		if parsed, errParse := uuid.Parse(uaIDStr); errParse == nil {
			userApartmentID = &parsed
		}
	}

	userRoleVal, _ := ctx.Get("role")
	userRole, _ := userRoleVal.(string)

	err = c.pollService.RevokeVote(pollID, userID, userApartmentID, userRole)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "vote_revoked_successfully", nil)
}

func (c *PollController) Delete(ctx *gin.Context) {
	pollIDStr := ctx.Param("poll_id")
	pollID, err := uuid.Parse(pollIDStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_poll_id", err)
		return
	}

	if err := c.pollService.DeletePoll(pollID); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "poll_deleted_successfully", nil)
}
