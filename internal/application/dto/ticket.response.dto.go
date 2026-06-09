package dto

import (
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"time"

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
	CreatedAt   time.Time  `json:"created_at"`
	Tags        []TagResponse
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
	CreatedAt   time.Time     `json:"created_at"`
}

type TicketBaseResponseWithComments struct {
	TicketBaseResponse
	Comments []CommentResponse `json:"comments"`
}

type TicketBaseResponseWithUser struct {
	TicketBaseResponse
	User      UserResponse      `json:"user"`
	Apartment ApartmentResponse `json:"apartment"`
}

type TicketBaseResponseWithAllRelations struct {
	TicketBaseResponse
	Comments  []CommentResponseWithUser `json:"comments"`
	User      UserResponse              `json:"user"`
	Apartment ApartmentResponse         `json:"apartment"`
}

type TicketBaseResponseWithCommentCount struct {
	TicketBaseResponse
	CommentCount int64 `json:"comment_count"`
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
		CreatedAt:   ticket.CreatedAt,
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
	if len(items) == 0 {
		return []TicketBaseResponseWithCommentCount{}
	}

	responses := make([]TicketBaseResponseWithCommentCount, 0, len(items))

	for i := range items {
		responses = append(responses, *MapTicketToResponseWithCount(&items[i].Ticket, items[i].CommentCount))
	}

	return responses
}
