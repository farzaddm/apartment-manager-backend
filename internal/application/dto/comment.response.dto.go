package dto

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type CreateCommentResponse struct {
	ID             uuid.UUID  `json:"id"`
	UserID         *uuid.UUID `json:"user_id"`
	TicketID       uuid.UUID  `json:"ticket_id"`
	Body           string     `json:"body"`
	CommittedOrder int        `json:"committed_order"`
}

type CommentResponse struct {
	ID             uuid.UUID
	UserID         *uuid.UUID
	TicketID       uuid.UUID
	Body           string
	CommittedOrder int
}

type CommentResponseWithUser struct {
	CommentResponse
	User UserResponse
}

type CommentResponseWithTicket struct {
	CommentResponse
	Ticket TicketBaseResponse
}

type CommentResponseWithAllRelations struct {
	CommentResponse
	User   UserResponse
	Ticket TicketBaseResponse
}

func MapCommentToResponse(comment *entity.Comment) *CommentResponse {
	if comment == nil {
		return nil
	}

	return &CommentResponse{
		ID:             comment.ID,
		UserID:         comment.UserID,
		TicketID:       comment.TicketID,
		Body:           comment.Body,
		CommittedOrder: comment.CommittedOrder,
	}
}

func MapCommentsToSliceResponse(comments []entity.Comment) []CommentResponse {
	if comments == nil {
		return nil
	}

	responses := make([]CommentResponse, 0, len(comments))

	for i := range comments {
		responses = append(responses, *MapCommentToResponse(&comments[i]))
	}

	return responses
}
