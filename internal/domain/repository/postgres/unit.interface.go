package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type UnitInterface interface {
	GetByID(ctx context.Context, id string) (*entity.Unit, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Unit, error)
}
