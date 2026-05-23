package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{db: db}
}

func (r *UnitRepository) GetByID(ctx context.Context, id string) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&unit).Error
	if err != nil {
		return nil, err
	}

	return &unit, nil
}

func (r *UnitRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).First(&unit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}
