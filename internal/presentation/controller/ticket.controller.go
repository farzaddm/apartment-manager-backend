package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	service_error "apartment-manager-backend/internal/application/service/error"
	"fmt"
	"log"

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

	data, err := c.ticketService.Create(ctx, tokenKeys, &req)
	if err != nil {
		switch err {
		case service_error.ErrTicketUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "ticket_unauthorized_access", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
			return
		}
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
		filter.UserUUID = nil
	} else {
		uid, err := uuid.Parse(b)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_ticket_id", err)
			return
		}
		filter.UserUUID = &uid
	}
	b = ctx.Query("status")
	if b == "" {
		filter.Status = nil
	}
	b = ctx.Query("category")
	if b == "" {
		filter.Category = nil
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

	tickets, err := c.ticketService.List(ctx, tokenKeys, filter)
	if err != nil {

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

	var data *dto.TicketBaseResponse
	if data, err = c.ticketService.Update(ctx, tokenKeys, id, req); err != nil {
		switch err {

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrTicketUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "ticket_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
			return
		}

	}

	response.Success(ctx, http.StatusOK, "ticket_updated_successfully", data)
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

	var data *dto.TicketBaseResponse
	if data, err = c.ticketService.UpdateStatus(ctx, tokenKeys, id, req.Status); err != nil {
		switch err {

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrTicketUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "ticket_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
			return
		}

	}

	response.Success(ctx, http.StatusOK, "ticket_status_updated_successfully", data)
}

func (c *TicketController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_ticket_id", err)
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

	if err := c.ticketService.Delete(ctx, tokenKeys, id); err != nil {
		switch err {

		case service_error.ErrTicketUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "ticket_unauthorized_access", err)
			return

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
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
	ticket, err := c.ticketService.GetByID(ctx, tokenKeys, id)

	if err != nil {
		switch err {

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrTicketUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "ticket_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
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

	ticket, err := c.ticketService.GetByIDWithAllRelations(ctx, tokenKeys, id)

	if err != nil {
		switch err {

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrTicketUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "ticket_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
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

	data, err := c.commentService.Create(ctx, tokenKeys, ticketID, &req)
	if err != nil {
		switch err {

		case service_error.ErrTicketNotFound:
			response.Error(ctx, http.StatusNotFound, "ticket_not_found", err)
			return

		case service_error.ErrCommentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "forbidden_comment_access", err)
			return

		case service_error.ErrCommentNotFound:
			response.Error(ctx, http.StatusNotFound, "comment_order_not_found", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "create_comment_failed", err)
			return
		}
	}

	response.Success(ctx, http.StatusCreated, "comment_created_successfully", data)
}
