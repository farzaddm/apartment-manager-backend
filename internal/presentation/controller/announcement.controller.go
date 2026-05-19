package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnnouncementController struct {
	announcementService service.AnnouncementService
}

func NewAnnouncementController(as service.AnnouncementService) *AnnouncementController {
	return &AnnouncementController{announcementService: as}
}

func (c *AnnouncementController) Create(ctx *gin.Context) {
	aptID, err := uuid.Parse(ctx.Param("apartment_id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	var req dto.CreateAnnouncementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	res, err := c.announcementService.Create(ctx.Request.Context(), aptID, req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusCreated, "announcement_created_successfully", res)
}

func (c *AnnouncementController) Get(ctx *gin.Context) {
	aptID, err := uuid.Parse(ctx.Param("apartment_id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_announcement_id", err)
		return
	}

	res, err := c.announcementService.GetByID(ctx.Request.Context(), id, aptID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "announcement_fetch_successfully", res)
}

func (c *AnnouncementController) Update(ctx *gin.Context) {
	aptID, err := uuid.Parse(ctx.Param("apartment_id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_announcement_id", err)
		return
	}

	var req dto.UpdateAnnouncementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	res, err := c.announcementService.Update(ctx.Request.Context(), id, aptID, req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "announcement_update_successfully", res)
}

func (c *AnnouncementController) Delete(ctx *gin.Context) {
	aptID, err := uuid.Parse(ctx.Param("apartment_id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_announcement_id", err)
		return
	}

	if err := c.announcementService.Delete(ctx.Request.Context(), id, aptID); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "announcement_delete_successfully", nil)
}
