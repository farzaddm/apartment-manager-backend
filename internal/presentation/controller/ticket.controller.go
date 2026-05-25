package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	service_error "apartment-manager-backend/internal/application/service/error"
	"fmt"

	"apartment-manager-backend/pkg/response"
	"apartment-manager-backend/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketController struct {
	ticketService  service.TicketService
	commentService service.CommentService
}

func NewTicketController(ticketService service.TicketService, commentService service.CommentService) *TicketController {
	return &TicketController{
		ticketService:  ticketService,
		commentService: commentService,
	}
}

func (c *TicketController) Create(ctx *gin.Context) {
	var req dto.CreateTicketRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errList := validator.ParseValidationErrors(err)

		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "invalid_request_body",
			Errors:     errList,
		}
		res.SendResponse(ctx)
		return
	}

	ticket := &dto.CreateTicketRequest{
		Title:         req.Title,
		Description:   req.Description,
		Body:          req.Body,
		Category:      req.Category,
		Accessability: req.Accessability,
	}

	cr_ticket, err := c.ticketService.Create(ctx, ticket)
	if err != nil {
		switch err {

		case service_error.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access", err)
			return

		case service_error.ErrCommonParseStrToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_uuid_format", err)
			return

		case service_error.ErrTicketUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "ticket_access_denied", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_create_ticket", err)
			return
		}
	}
	data := dto.CreateTicketResponse{
		ID:          cr_ticket.ID,
		UserID:      cr_ticket.UserID,
		Title:       cr_ticket.Title,
		Description: cr_ticket.Description,
		Body:        cr_ticket.Body,
		Category:    cr_ticket.Category,
		Status:      cr_ticket.Status,
	}
	response.Success(ctx, http.StatusCreated, "ticket_created_successfully", data)
}

func (c *TicketController) List(ctx *gin.Context) {
	var filter dto.TicketFilterRequest

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_query_params", err)
		return
	}

	var b string

	b = ctx.Query("user_id")
	if b == "" {
		filter.UserID = nil
	}
	b = ctx.Query("status")
	if b == "" {
		filter.Status = nil
	}
	b = ctx.Query("category")
	if b == "" {
		filter.Category = nil
	}

	tickets, err := c.ticketService.List(ctx, filter)
	if err != nil {
		switch err {
		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrUserIDNotFoundInContext,
			service_error.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access", err)
			return

		case service_error.ErrCommonParseStrToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_uuid_format", err)
			return

		case service_error.ErrTicketIsPrivate:
			response.Error(ctx, http.StatusForbidden, "ticket_is_private_and_access_denied", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_fetch_tickets", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "tickets_fetched_successfully", tickets)
}

func (c *TicketController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_ticket_id", err)
		return
	}

	var req dto.UpdateTicketRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errList := validator.ParseValidationErrors(err)

		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "invalid_request_body",
			Errors:     errList,
		}
		res.SendResponse(ctx)
		return
	}

	if err := c.ticketService.Update(ctx, id, req); err != nil {
		switch err {

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access", err)
			return

		case service_error.ErrCommonParseStrToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_uuid_format", err)
			return

		case service_error.ErrTicketUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "ticket_access_denied", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_update_ticket", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "ticket_updated_successfully", nil)
}

func (c *TicketController) UpdateTicketStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_ticket_id", err)
		return
	}

	var req dto.UpdateTicketStatusRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errList := validator.ParseValidationErrors(err)

		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "invalid_request_body",
			Errors:     errList,
		}
		res.SendResponse(ctx)
		return
	}

	if err := c.ticketService.UpdateStatus(ctx, id, req.Status); err != nil {
		switch err {
		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_update_ticket_status", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "ticket_status_updated_successfully", nil)
}

func (c *TicketController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_ticket_id", err)
		return
	}

	if err := c.ticketService.Delete(ctx, id); err != nil {
		switch err {

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrUserIDNotFoundInContext,
			service_error.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access", err)
			return

		case service_error.ErrCommonParseStrToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_uuid_format", err)
			return

		case service_error.ErrTicketUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "ticket_access_denied", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_delete_ticket", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "ticket_deleted_successfully", nil)
}

func (c *TicketController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_ticket_id", err)
		return
	}

	ticket, err := c.ticketService.GetByID(ctx, id)

	if err != nil {
		switch err {
		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_missing_in_context", err)
			return

		case service_error.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_missing_in_context", err)
			return

		case service_error.ErrTicketIsPrivate:
			response.Error(ctx, http.StatusForbidden, "ticket_is_private_and_access_denied", err)
			return

		case service_error.ErrCommonParseStrToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_fetch_ticket", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "ticket_fetched_successfully", ticket)
}

func (c *TicketController) GetFully(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_ticket_id", err)
		return
	}

	ticket, err := c.ticketService.GetByIDWithAllRelations(ctx, id)

	if err != nil {
		switch err {
		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrUserIDNotFoundInContext,
			service_error.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access", err)
			return

		case service_error.ErrTicketIsPrivate:
			response.Error(ctx, http.StatusForbidden, "ticket_is_private_and_access_denied", err)
			return

		case service_error.ErrCommonParseStrToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_uuid_format", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_fetch_ticket", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "ticket_fetched_successfully", ticket)
}

func (c *TicketController) GetUserTickets(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "user_id_not_found", nil)
		return
	}

	tickets, err := c.ticketService.GetUserTickets(ctx.Request.Context(), fmt.Sprintf("%v", userID))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_fetch_tickets", err)
		return
	}

	response.Success(ctx, http.StatusOK, "tickets_fetched_successfully", tickets)
}

func (c *TicketController) CreateComment(ctx *gin.Context) {
	ticketID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_ticket_id", err)
		return
	}

	var req dto.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errList := validator.ParseValidationErrors(err)

		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "invalid_request_body",
			Errors:     errList,
		}
		res.SendResponse(ctx)
		return
	}

	comm, err := c.commentService.Create(ctx, ticketID, &req)
	if err != nil {
		switch err {

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
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

		case service_error.ErrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_not_found", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_create_comment", err)
			return
		}
	}
	data := dto.CreateCommentResponse{
		ID:             comm.ID,
		UserID:         comm.UserID,
		TicketID:       comm.TicketID,
		Body:           comm.Body,
		CommittedOrder: comm.CommittedOrder,
	}
	response.Success(ctx, http.StatusCreated, "comment_created_successfully", data)
}
