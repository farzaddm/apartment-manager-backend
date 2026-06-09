package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/pkg/response"
	"apartment-manager-backend/pkg/validator"
	"log"
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

	tokenKeys, err := dto.NewTokenKeys(ctx)
	if err != nil {
		switch err {
		case dto.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_not_found_in_context", err)
			return
		case dto.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_not_found_in_context", err)
			return
		case dto.ErrApartmentIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "apartment_id_not_found_in_context", err)
			return
		case dto.ErrUserIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return
		case dto.ErrApartmentIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id_format", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "unexpected_error", err)
			log.Println(err)
			return
		}
	}

	comment, err := h.commentService.GetByID(ctx, tokenKeys, id)
	if err != nil {
		switch err {

		case service_error.ErrTicketOfCommentOrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_or_ticket_not_found", err)
			return

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_not_found", err)
			return

		case service_error.ErrCommentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "forbidden_comment_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "get_comment_failed", err)
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
		errList := validator.ParseValidationErrors(err)

		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "invalid_request_body",
			Errors:     append(errList, err.Error()),
		}
		res.SendResponse(ctx)
		return
	}

	tokenKeys, err := dto.NewTokenKeys(ctx)
	if err != nil {
		switch err {
		case dto.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_not_found_in_context", err)
			return
		case dto.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_not_found_in_context", err)
			return
		case dto.ErrApartmentIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "apartment_id_not_found_in_context", err)
			return
		case dto.ErrUserIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return
		case dto.ErrApartmentIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id_format", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "unexpected_error", err)
			log.Println(err)
			return
		}
	}

	data, err := h.commentService.Update(ctx, tokenKeys, id, &req)
	if err != nil {
		switch err {

		case service_error.ErrTicketOfCommentOrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_or_ticket_not_found", err)
			return

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_not_found", err)
			return

		case service_error.ErrCommentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "forbidden_comment_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "update_comment_failed", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "comment_updated_successfully", data)
}

func (h *CommentController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_comment_id", err)
		return
	}

	tokenKeys, err := dto.NewTokenKeys(ctx)
	if err != nil {
		switch err {
		case dto.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_not_found_in_context", err)
			return
		case dto.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_not_found_in_context", err)
			return
		case dto.ErrApartmentIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "apartment_id_not_found_in_context", err)
			return
		case dto.ErrUserIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return
		case dto.ErrApartmentIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id_format", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "unexpected_error", err)
			log.Println(err)
			return
		}
	}

	if err := h.commentService.Delete(ctx, tokenKeys, id); err != nil {
		switch err {

		case service_error.ErrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_not_found", err)
			return

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrCommentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "forbidden_comment_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "delete_comment_failed", err)
			return

		}
	}

	response.Success(ctx, http.StatusOK, "comment_deleted_successfully", nil)
}
