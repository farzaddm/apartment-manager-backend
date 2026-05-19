package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/internal/domain/entity"
	"fmt"

	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketController struct {
	ticketService service.TicketService
}

func NewTicketController(ticketService service.TicketService) *TicketController {
	return &TicketController{
		ticketService: ticketService,
	}
}

func (c *TicketController) Create(ctx *gin.Context) {
	var req dto.CreateTicketRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	ticket := &entity.Ticket{
		UserID:        req.UserID, // must be included
		Title:         req.Title,
		Description:   req.Description,
		Body:          req.Body,
		Category:      req.Category,
		Accessability: req.Accessability,
	}

	baseUserID, exists := ctx.Get("user_id") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_userid"))
		return
	}

	//TODO:check other errors
	if err := c.ticketService.Create(ctx, baseUserID.(uuid.UUID), ticket); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_create_ticket", err)
		return
	}

	response.Success(ctx, http.StatusCreated, "ticket_created_successfully", ticket)
}

func (c *TicketController) List(ctx *gin.Context) {
	var filter dto.TicketFilterRequest

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_query_params", err)
		return
	}

	var b bool

	_, b = ctx.Get("user_id")
	if b == false {
		filter.UserID = nil
	}
	_, b = ctx.Get("status")
	if b == false {
		filter.Status = nil
	}
	_, b = ctx.Get("category")
	if b == false {
		filter.Category = nil
	}

	baseUserID, exists := ctx.Get("user_id") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_userid"))
		return
	}
	role, exists := ctx.Get("role") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_user_role"))
		return
	}

	tickets, err := c.ticketService.List(ctx, baseUserID.(uuid.UUID), filter, role.(entity.UserRole))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_fetch_tickets", err)
		return
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
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	baseUserID, exists := ctx.Get("user_id") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_userid"))
		return
	}

	if err := c.ticketService.Update(ctx, baseUserID.(uuid.UUID), id, req); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_update_ticket", err)
		return
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
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	if err := c.ticketService.UpdateStatus(ctx, id, req.Status); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_update_ticket_status", err)
		return
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

	baseUserID, exists := ctx.Get("user_id") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_userid"))
		return
	}
	role, exists := ctx.Get("role") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_user_role"))
		return
	}

	if err := c.ticketService.Delete(ctx, baseUserID.(uuid.UUID), id, role.(entity.UserRole)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_delete_ticket", err)
		return
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

	baseUserID, exists := ctx.Get("user_id") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_userid"))
		return
	}
	role, exists := ctx.Get("role") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_user_role"))
		return
	}

	ticket, err := c.ticketService.GetByID(ctx, baseUserID.(uuid.UUID), id, role.(entity.UserRole))
	//TODO:check other errors
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_fetch_ticket", err)
		return
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

	baseUserID, exists := ctx.Get("user_id") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_userid"))
		return
	}
	role, exists := ctx.Get("role") // IT MUST BE EXIST!
	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "not_expected_authorization_action", fmt.Errorf("not_expected_authorization_action_user_role"))
		return
	}

	ticket, err := c.ticketService.GetByIDWithAllRelations(ctx, baseUserID.(uuid.UUID), id, role.(entity.UserRole))
	//TODO:check other errors
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_fetch_ticket", err)
		return
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
