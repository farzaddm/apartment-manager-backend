package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
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

	ticket := &dto.CreateTicketRequest{
		UserID:        req.UserID, // must be included
		Title:         req.Title,
		Description:   req.Description,
		Body:          req.Body,
		Category:      req.Category,
		Accessability: req.Accessability,
	}

	//TODO:check other errors
	if err := c.ticketService.Create(ctx, ticket); err != nil {
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

	var b string
	//TODO : check q params! conflict to normal key-values in ctx
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

	if err := c.ticketService.Update(ctx, id, req); err != nil {
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

	if err := c.ticketService.Delete(ctx, id); err != nil {
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

	ticket, err := c.ticketService.GetByID(ctx, id)
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

	ticket, err := c.ticketService.GetByIDWithAllRelations(ctx, id)
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
