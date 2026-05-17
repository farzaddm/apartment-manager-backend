package dto

import "apartment-manager-backend/internal/domain/entity"

type UpdateTicketRequest struct {
	Status      entity.TicketStatus
	Category    entity.TicketCategory
	Title       string
	Description string
	Body        string
}
