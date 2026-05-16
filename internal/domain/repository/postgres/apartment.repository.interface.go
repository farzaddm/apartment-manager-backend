package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"
)

type ApartmentInterface interface {
	Create(ctx context.Context, apartment *entity.Apartment) error
	Update(ctx context.Context, apartment *entity.Apartment) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (*bool, error)

	GetByID(ctx context.Context, id string) (*entity.Apartment, error)

	GetByIDWithRelations(ctx context.Context, id string, relations ...string) (*entity.Apartment, error)

	GetWithUsers(ctx context.Context, id string) (*entity.Apartment, error)
	GetWithAnnouncements(ctx context.Context, id string) (*entity.Apartment, error)
	GetWithRules(ctx context.Context, id string) (*entity.Apartment, error)
	GetWithInviteCodes(ctx context.Context, id string) (*entity.Apartment, error)

	List(ctx context.Context) ([]entity.Apartment, error)

	ListWithRelations(ctx context.Context, relations ...string) ([]entity.Apartment, error)

	ListWithUsers(ctx context.Context) ([]entity.Apartment, error)
	ListWithAnnouncements(ctx context.Context) ([]entity.Apartment, error)
	ListWithRules(ctx context.Context) ([]entity.Apartment, error)
	ListWithInviteCodes(ctx context.Context) ([]entity.Apartment, error)
}
