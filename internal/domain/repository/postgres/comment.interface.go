package postgres

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type CommentInterface interface {
	Create(comment *entity.Comment) error

	GetByID(id uuid.UUID) (*entity.Comment, error)

	ListByTicketID(ticketID uuid.UUID) ([]entity.Comment, error)

	Update(comment *entity.Comment) error

	Delete(id uuid.UUID) error

	GetLastOrderByTicketID(ticketID uuid.UUID) (*int, error)
}
