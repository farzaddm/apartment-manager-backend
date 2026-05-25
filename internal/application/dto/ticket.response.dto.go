package dto

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
)

// import (
// 	"apartment-manager-backend/internal/domain/entity"

// 	"github.com/google/uuid"
// )

// type TicketResponse struct {
// 	ID          uuid.UUID `json:"id"`
// 	UserID      uuid.UUID `json:"user_id"`
// 	Title       string    `json:"title"`
// 	Description string    `json:"description"`
// 	Body        string    `json:"body"`
// 	Category    string    `json:"category"`
// 	Status      string    `json:"status"`
// 	Tags        []string  `json:"tags,omitempty"` // Returns an array of clean Tag IDs string
// }

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

// type TicketBaseResponseWithComments struct {
// 	TicketBaseResponse
// 	Comments []dto.XComment
// }

// type TicketBaseResponseWithUser struct {
// 	TicketBaseResponse
// 	User UserResponse
// }

// type TicketBaseResponseWithAllRelations struct {
// 	TicketBaseResponse
// 	Comments []dto.XComment
// 	User     UserResponse
// }

func MapTicketToBaseResponse(ticket *entity.Ticket) *TicketBaseResponse {
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
