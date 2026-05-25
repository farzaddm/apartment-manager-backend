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
	ID             uuid.UUID  `json:"id"`
	UserID         *uuid.UUID `json:"user_id"`
	TicketID       uuid.UUID  `json:"ticket_id"`
	Body           string     `json:"body"`
	CommittedOrder int        `json:"committed_order"`
}

type CommentResponseWithUser struct {
	CommentResponse
	User UserResponse `json:"user"`
}

type CommentResponseWithTicket struct {
	CommentResponse
	Ticket TicketBaseResponse `json:"ticket"`
}

type CommentResponseWithAllRelations struct {
	CommentResponse
	User   UserResponse       `json:"user"`
	Ticket TicketBaseResponse `json:"ticket"`
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
	if len(comments) == 0 {
		return []CommentResponse{}
	}

	responses := make([]CommentResponse, 0, len(comments))

	for i := range comments {
		responses = append(responses, *MapCommentToResponse(&comments[i]))
	}

	return responses
}

func MapCommentToResponseWithUser(comment *entity.Comment) *CommentResponseWithUser {
	if comment == nil {
		return nil
	}

	return &CommentResponseWithUser{
		CommentResponse: *MapCommentToResponse(comment),
		User:            *MapUserToUserResponse(&comment.User),
	}
}

func MapCommentToResponseWithTicket(comment *entity.Comment) *CommentResponseWithTicket {
	if comment == nil {
		return nil
	}

	return &CommentResponseWithTicket{
		CommentResponse: *MapCommentToResponse(comment),
		Ticket:          *MapTicketToBaseResponse(&comment.Ticket),
	}
}

func MapCommentsToResponseWithUserSlice(comments []entity.Comment) []CommentResponseWithUser {
	if len(comments) == 0 {
		return []CommentResponseWithUser{}
	}

	responses := make([]CommentResponseWithUser, 0, len(comments))
	for i := range comments {
		mapped := MapCommentToResponseWithUser(&comments[i])
		if mapped != nil {
			responses = append(responses, *mapped)
		}
	}

	return responses
}
