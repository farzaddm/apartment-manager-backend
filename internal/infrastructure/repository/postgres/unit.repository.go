package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{db: db}
}

func (r *UnitRepository) GetByID(ctx context.Context, id string) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&unit).Error
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

func (r *UnitRepository) Create(ctx context.Context, unit *entity.Unit) error {
	return r.db.WithContext(ctx).Create(unit).Error
}

func (r *UnitRepository) Update(ctx context.Context, unit *entity.Unit) (*entity.Unit, error) {
	result := r.db.WithContext(ctx).
		Model(unit).
		Where("id = ?", unit.ID).
		Clauses(clause.Returning{}).
		Updates(unit)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return unit, nil
}

func (r *UnitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Delete(&entity.Unit{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UnitRepository) PopUser(ctx context.Context, id uuid.UUID) (*entity.Unit, error) {
	var unit entity.Unit

	result := r.db.WithContext(ctx).
		Model(&unit).
		Where("id = ?", id).
		Update("user_id", nil)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	err := r.db.WithContext(ctx).First(&unit, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &unit, nil
}

func (r *UnitRepository) PushUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Unit, error) {
	var unit entity.Unit

	result := r.db.WithContext(ctx).
		Model(&unit).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Update("user_id", userID)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &unit, nil
}
