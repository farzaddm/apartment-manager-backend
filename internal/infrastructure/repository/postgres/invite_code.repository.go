package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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

func (r *InviteCodeRepository) GetByCode(ctx context.Context, code string) (*entity.InviteCode, error) {
	var invite entity.InviteCode
	err := r.db.WithContext(ctx).Where("code = ?", code).Preload("Unit").First(&invite).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if code doesn't exist
		}
		return nil, err
	}
	return &invite, nil
}

func (r *InviteCodeRepository) AssignUserToUnitAndApartment(ctx context.Context, userID string, apartmentID string, unitID string) error {
	userUUID, _ := uuid.Parse(userID)
	aptUUID, _ := uuid.Parse(apartmentID)
	unitUUID, _ := uuid.Parse(unitID)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.User{}).Where("id = ?", userUUID).Update("apartment_id", aptUUID).Error
		if err != nil {
			return err
		}

		err = tx.Model(&entity.Unit{}).Where("id = ?", unitUUID).Update("user_id", userUUID).Error
		if err != nil {
			return err
		}

		return nil
	})
}
