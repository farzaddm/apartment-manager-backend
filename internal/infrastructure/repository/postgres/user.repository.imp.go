package postgres

import (
	"context"
	"errors"

	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainRepo.UserInterface {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("phone = ?", phone).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) ExistEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) ExistPhone(ctx context.Context, phone string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
