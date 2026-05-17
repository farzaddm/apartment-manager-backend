package dto

import "apartment-manager-backend/internal/domain/entity"

type UpdateTicketRequest struct {
	Category    entity.TicketCategory
	Title       string
	Description string
	Body        string
}
