package dto

import (
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"

	"github.com/google/uuid"
)

type TicketResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
	Tags        []string  `json:"tags,omitempty"` // Returns an array of clean Tag IDs string
}

type CreateTicketResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      *uuid.UUID `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Body        string     `json:"body"`
	Category    string     `json:"category"`
	Status      string     `json:"status"`
}

// // TODO : Coupling TicketResponse!!!!!! Remove one of them!!!!!
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

type TicketBaseResponseWithComments struct {
	TicketBaseResponse
	Comments []CommentResponse
}

type TicketBaseResponseWithUser struct {
	TicketBaseResponse
	User UserResponse
}

type TicketBaseResponseWithAllRelations struct {
	TicketBaseResponse
	Comments []CommentResponse
	User     UserResponse
}

type TicketBaseResponseWithCommentCount struct {
	TicketBaseResponse
	CommentCount int64
}

func MapTicketToBaseResponse(ticket *entity.Ticket) *TicketBaseResponse {
	if ticket == nil {
		return nil
	}

	var tags = make([]TagResponse, 0)

	for i := range ticket.Tags {
		tags = append(tags, *MapTagToResponse(&ticket.Tags[i].Tag))
	}

	return &TicketBaseResponse{
		ID:          ticket.ID,
		UserID:      ticket.UserID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Body:        ticket.Body,
		Category:    string(ticket.Category),
		Status:      string(ticket.Status),
		Tags:        tags, //TODO : Fix Ticket-Tags Relation
	}
}

func MapTicketToResponseWithCount(ticket *entity.Ticket, count int64) *TicketBaseResponseWithCommentCount {
	if ticket == nil {
		return nil
	}

	return &TicketBaseResponseWithCommentCount{
		TicketBaseResponse: *MapTicketToBaseResponse(ticket),
		CommentCount:       count,
	}
}

func MapTicketsToSliceResponseWithCount(items []domainRepo.TicketWithCommentCount) []TicketBaseResponseWithCommentCount {
	if items == nil {
		return nil
	}

	responses := make([]TicketBaseResponseWithCommentCount, 0, len(items))

	for i := range items {
		responses = append(responses, *MapTicketToResponseWithCount(&items[i].Ticket, items[i].CommentCount))
	}

	return responses
}
