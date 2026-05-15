package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"
)

type UserInterface interface {
	Create(ctx context.Context, user *entity.User) error
	ExistEmail(ctx context.Context, email string) (bool, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
}
