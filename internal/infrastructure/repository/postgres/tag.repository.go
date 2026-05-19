package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) postgres.TagInterface {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *TagRepository) FindAll(ctx context.Context) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Tag{}, "id = ?", id).Error
}
