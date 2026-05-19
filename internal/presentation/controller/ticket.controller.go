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
		response.Error(ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	ticket := &entity.Ticket{
		UserID:      req.UserID, // must be included
		Title:       req.Title,
		Description: req.Description,
		Body:        req.Body,
		Category:    req.Category,
	}

	if err := c.ticketService.Create(ctx, ticket); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create ticket", err)
		return
	}

	response.Success(ctx, http.StatusCreated, "ticket created successfully", ticket)
}

func (c *TicketController) List(ctx *gin.Context) {
	var filter dto.TicketFilterRequest

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid query params", err)
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

	tickets, err := c.ticketService.List(ctx, filter)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch tickets", err)
		return
	}

	response.Success(ctx, http.StatusOK, "tickets fetched successfully", tickets)
}

func (c *TicketController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ticket id", err)
		return
	}

	var req dto.UpdateTicketRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := c.ticketService.Update(ctx, id, req); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update ticket", err)
		return
	}

	response.Success(ctx, http.StatusOK, "ticket updated successfully", nil)
}

func (c *TicketController) UpdateTicketStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ticket id", err)
		return
	}

	var req dto.UpdateTicketStatusRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := c.ticketService.UpdateStatus(ctx, id, req.Status); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update ticket status", err)
		return
	}

	response.Success(ctx, http.StatusOK, "ticket status updated successfully", nil)
}

func (c *TicketController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ticket id", err)
		return
	}

	if err := c.ticketService.Delete(ctx, id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete ticket", err)
		return
	}

	response.Success(ctx, http.StatusOK, "ticket deleted successfully", nil)
}

func (c *TicketController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ticket id", err)
		return
	}

	ticket, err := c.ticketService.GetByID(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch ticket", err)
		return
	}

	response.Success(ctx, http.StatusOK, "ticket fetched successfully", ticket)
}

func (c *TicketController) GetFully(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ticket id", err)
		return
	}

	ticket, err := c.ticketService.GetByIDWithAllRelations(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch ticket", err)
		return
	}

	response.Success(ctx, http.StatusOK, "ticket fetched successfully", ticket)
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
