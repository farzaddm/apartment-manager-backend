package dto

import "github.com/google/uuid"

type CreateCommentRequest struct {
	TicketID uuid.UUID `json:"ticket_id" binding:"required"`
	Body     string    `json:"body" binding:"required"`
}

type UpdateCommentRequest struct {
	Body string `json:"body"`
}
