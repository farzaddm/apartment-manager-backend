package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type UnitInterface interface {
	Create(ctx context.Context, unit *entity.Unit) error
	Update(ctx context.Context, unit *entity.Unit) (*entity.Unit, error)
	Delete(ctx context.Context, id uuid.UUID) error
	PopUser(ctx context.Context, id uuid.UUID) (*entity.Unit, error)
	GetByID(ctx context.Context, id string) (*entity.Unit, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Unit, error)
	PushUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Unit, error)
}
