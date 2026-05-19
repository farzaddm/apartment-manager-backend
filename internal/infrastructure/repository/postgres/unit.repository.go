package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"

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
