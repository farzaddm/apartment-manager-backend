package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type TagInterface interface {
	Create(ctx context.Context, tag *entity.Tag) error
	FindAll(ctx context.Context) ([]entity.Tag, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
