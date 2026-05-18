package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TagController struct {
	tagService service.TagService
}

func NewTagController(tagService service.TagService) *TagController {
	return &TagController{tagService: tagService}
}

func (c *TagController) Create(ctx *gin.Context) {
	var req dto.CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	res, err := c.tagService.Create(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusCreated, "tag_created_successfully", res)
}

func (c *TagController) List(ctx *gin.Context) {
	res, err := c.tagService.GetAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "tag_list_fetch_successfully", res)
}

func (c *TagController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalide_tag_id_format", err)
		return
	}

	if err := c.tagService.Delete(ctx.Request.Context(), id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "tag_deleted_successfully", nil)
}
