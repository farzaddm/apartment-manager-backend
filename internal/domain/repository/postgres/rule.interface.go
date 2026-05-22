package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type RuleRepository interface {
	Create(ctx context.Context, rule *entity.Rule) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Rule, error)
	GetByApartmentID(ctx context.Context, apartmentID uuid.UUID) ([]entity.Rule, error)
	GetByApartmentAndCategory(ctx context.Context, apartmentID uuid.UUID, category entity.RuleCategory) ([]entity.Rule, error)

	Update(ctx context.Context, rule *entity.Rule) error
	Delete(ctx context.Context, id uuid.UUID) error
}
