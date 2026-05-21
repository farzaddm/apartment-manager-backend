package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	service_error "apartment-manager-backend/internal/application/service/error"
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
		switch err {
		case service_error.ErrTicketOfCommentOrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_or_comment_not_found", err)
			return

		case service_error.ErrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_not_found", err)
			return

		case service_error.ErrUserIDNotFoundInContext,
			service_error.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access", err)
			return

		case service_error.ErrCommonParseStrToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_uuid_format", err)
			return

		case service_error.ErrCommentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "comment_access_denied", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_fetch_comment", err)
			return
		}
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
		switch err {

		case service_error.ErrTicketOfCommentOrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_or_comment_not_found", err)
			return

		case service_error.ErrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_not_found", err)
			return

		case service_error.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access", err)
			return

		case service_error.ErrCommonParseStrToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_uuid_format", err)
			return

		case service_error.ErrCommentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "comment_access_denied", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_update_comment", err)
			return
		}
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
		switch err {

		case service_error.ErrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_not_found", err)
			return

		case service_error.ErrUserIDNotFoundInContext,
			service_error.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access", err)
			return

		case service_error.ErrCommonParseStrToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_uuid_format", err)
			return

		case service_error.ErrCommentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "comment_access_denied", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_delete_comment", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "comment_deleted_successfully", nil)
}
