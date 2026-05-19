package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"
)

type UnitInterface interface {
	GetByID(ctx context.Context, id string) (*entity.Unit, error)
}
