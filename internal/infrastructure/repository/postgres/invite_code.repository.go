package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type InviteCodeRepository struct {
	db *gorm.DB
}

func NewInviteCodeRepository(db *gorm.DB) *InviteCodeRepository {
	return &InviteCodeRepository{db: db}
}

func (r *InviteCodeRepository) Create(ctx context.Context, inviteCode *entity.InviteCode) error {
	return r.db.WithContext(ctx).Create(inviteCode).Error
}

func (r *InviteCodeRepository) GetActiveInviteCode(ctx context.Context, apartmentID string, unitID string) (*entity.InviteCode, error) {
	var invite entity.InviteCode

	err := r.db.WithContext(ctx).
		Where("apartment_id = ? AND unit_id = ? AND expires_at > ?", apartmentID, unitID, time.Now()).
		First(&invite).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active code found, which is what we want!
		}
		return nil, err // A real DB connection error happened
	}

	return &invite, nil
}
