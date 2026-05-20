package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentController struct {
	commentService service.CommentService
}

func NewCommentController(commentService service.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

func (h *CommentController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_comment_id", err)
		return
	}

	comment, err := h.commentService.GetByID(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "comment_not_found", err)
		return
	}

	response.Success(ctx, http.StatusOK, "comment_fetched_successfully", comment)
}

func (h *CommentController) UpdateBody(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_comment_id", err)
		return
	}

	var req dto.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	err = h.commentService.Update(ctx, id, &req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_update_comment", err)
		return
	}

	response.Success(ctx, http.StatusOK, "comment_updated_successfully", nil)
}

func (h *CommentController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_comment_id", err)
		return
	}

	if err := h.commentService.Delete(ctx, id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_delete_comment", err)
		return
	}

	response.Success(ctx, http.StatusOK, "comment_deleted_successfully", nil)
}
