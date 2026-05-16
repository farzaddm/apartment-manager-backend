package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type TicketInterface interface {
	Create(ctx context.Context, ticket *entity.Ticket) error

	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)

	GetByIDWithAllRelations(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)

	List(ctx context.Context, filter TicketFilter) ([]TicketWithCommentCount, error)

	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TicketStatus) error
	UpdateCategory(ctx context.Context, id uuid.UUID, category entity.TicketCategory) error
	UpdateContent(ctx context.Context, id uuid.UUID, title, description, body string) error
}

type TicketFilter struct {
	UserID   *uuid.UUID
	Status   *entity.TicketStatus
	Category *entity.TicketCategory

	Limit  *int
	Offset *int
}

type TicketWithCommentCount struct { 
	entity.Ticket
	CommentCount int64
}
