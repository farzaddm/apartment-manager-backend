package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ruleRepository struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) postgres.RuleRepository {
	return &ruleRepository{db: db}
}

func (r *ruleRepository) Create(ctx context.Context, rule *entity.Rule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *ruleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Rule, error) {
	var rule entity.Rule
	err := r.db.WithContext(ctx).First(&rule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *ruleRepository) GetByApartmentID(ctx context.Context, apartmentID uuid.UUID) ([]entity.Rule, error) {
	var rules []entity.Rule
	err := r.db.WithContext(ctx).Where("apartment_id = ? AND deleted_at IS NULL", apartmentID).Find(&rules).Error
	return rules, err
}

func (r *ruleRepository) GetByApartmentAndCategory(ctx context.Context, apartmentID uuid.UUID, category entity.RuleCategory) ([]entity.Rule, error) {
	var rules []entity.Rule
	err := r.db.WithContext(ctx).Where("apartment_id = ? AND category = ? AND deleted_at IS NULL", apartmentID, category).Find(&rules).Error
	return rules, err
}

func (r *ruleRepository) Update(ctx context.Context, rule *entity.Rule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

func (r *ruleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Rule{}, "id = ?", id).Error
}
