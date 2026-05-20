package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type ApartmentInterface interface {
	Create(ctx context.Context, apartment *entity.Apartment) error
	Update(ctx context.Context, id uuid.UUID, apartment *entity.Apartment) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (*bool, error)

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)

	GetByIDWithRelations(ctx context.Context, id uuid.UUID, relations ...string) (*entity.Apartment, error)

	GetWithUsers(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)
	GetWithAnnouncements(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)
	GetWithRules(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)
	GetWithInviteCodes(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)

	List(ctx context.Context) ([]entity.Apartment, error)

	ListWithRelations(ctx context.Context, relations ...string) ([]entity.Apartment, error)

	ListWithUsers(ctx context.Context) ([]entity.Apartment, error)
	ListWithAnnouncements(ctx context.Context) ([]entity.Apartment, error)
	ListWithRules(ctx context.Context) ([]entity.Apartment, error)
	ListWithInviteCodes(ctx context.Context) ([]entity.Apartment, error)

	GetApartmentManagerID(ctx context.Context, apartmentID string) (string, error)
}
