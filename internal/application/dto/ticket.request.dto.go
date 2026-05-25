package dto

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	Title         string                     `json:"title" binding:"required"`
	Description   string                     `json:"description"`
	Body          string                     `json:"body"  binding:"required"`
	Category      entity.TicketCategory      `json:"category" binding:"required"`
	Accessability entity.TicketAccessability `json:"accessability" binding:"required"`
}

// TODO : Coupling TicketResponse!!!!!! Remove one of them!!!!!
type TicketBaseResponse struct {
	ID          uuid.UUID     `json:"id"`
	UserID      *uuid.UUID    `json:"user_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Body        string        `json:"body"`
	Category    string        `json:"category"`
	Status      string        `json:"status"`
	Tags        []TagResponse `json:"tags"`
}

type UpdateTicketRequest struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Body        string                `json:"body"`
	Category    entity.TicketCategory `json:"category"`
}

type UpdateTicketStatusRequest struct {
	Status entity.TicketStatus `json:"status" binding:"required"`
}

type TicketFilterRequest struct {
	UserID   *uuid.UUID             `form:"user_id"`
	Status   *entity.TicketStatus   `form:"status"`
	Category *entity.TicketCategory `form:"category"`

	Page  int `form:"page"  binding:"required"`
	Limit int `form:"limit"  binding:"required"`
}

func MapTicketTotResponse(ticket *entity.Ticket) *TicketBaseResponse {
	if ticket == nil {
		return nil
	}
	return &TicketBaseResponse{
		ID:          ticket.ID,
		UserID:      ticket.UserID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Body:        ticket.Body,
		Category:    string(ticket.Category),
		Status:      string(ticket.Status),
		Tags:        nil, //TODO : Fix Ticket-Tags Relation
	}
}
