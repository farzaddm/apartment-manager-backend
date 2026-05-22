package postgres

import (
	"apartment-manager-backend/internal/application/dto/user"
	"apartment-manager-backend/internal/domain/entity"
	"context"
)

type UserInterface interface {
	Create(ctx context.Context, user *entity.User) error
	GetById(ctx context.Context, id string) (*entity.User, error)
	Update(ctx context.Context, user userdto.UpdateProfileRequest, id string) error
	ChangePassword(ctx context.Context, hashedPassword string, id string) error
	Delete(ctx context.Context, id string) error

	ExistEmail(ctx context.Context, email string) (bool, error)

	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	UpdateProfileImage(ctx context.Context, userID string, imagePath string) error
}
