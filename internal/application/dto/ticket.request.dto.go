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



