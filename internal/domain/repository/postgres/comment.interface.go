package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type CommentInterface interface {
	Create(ctx context.Context, comment *entity.Comment) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error)
	GetWithUser(ctx context.Context, id uuid.UUID) (*entity.Comment, error)
	GetWithTicket(ctx context.Context, id uuid.UUID) (*entity.Comment, error)
	GetWithAllRelations(ctx context.Context, id uuid.UUID) (*entity.Comment, error)

	ListByTicketID(ctx context.Context, ticketID uuid.UUID) ([]entity.Comment, error)

	Update(ctx context.Context, id uuid.UUID, comment *entity.Comment) (*entity.Comment, error)

	Delete(ctx context.Context, id uuid.UUID) error

	GetLastOrderByTicketID(ctx context.Context, ticketID uuid.UUID) (*int, error)
}
