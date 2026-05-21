package dto

import "github.com/google/uuid"

type CreateCommentRequest struct {
	Body string `json:"body" binding:"required"`
}

type CreateCommentResponse struct {
	ID             uuid.UUID  `json:"id"`
	UserID         *uuid.UUID `json:"user_id"`
	TicketID       uuid.UUID  `json:"ticket_id"`
	Body           string     `json:"body"`
	CommittedOrder int        `json:"committed_order"`
}

type UpdateCommentRequest struct {
	Body string `json:"body"`
}
