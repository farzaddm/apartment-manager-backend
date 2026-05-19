package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) postgres.AnnouncementInterface {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) Create(ctx context.Context, announcement *entity.Announcement) error {
	return r.db.WithContext(ctx).Create(announcement).Error
}

func (r *announcementRepository) FindByIDAndApartment(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) (*entity.Announcement, error) {
	var announcement entity.Announcement
	err := r.db.WithContext(ctx).
		Preload("Tags.Tag").
		Where("id = ? AND apartment_id = ?", id, apartmentID).
		First(&announcement).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}

func (r *announcementRepository) Update(ctx context.Context, announcement *entity.Announcement) error {
	// Updates primitive values on the core table omitting automatic associative processing hooks
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: false}).Save(announcement).Error
}

func (r *announcementRepository) ReplaceTags(ctx context.Context, announcementID uuid.UUID, tags []entity.TicketAnnouncementTag) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Wipe previous join mappings for this announcement
		if err := tx.Where("announcement_id = ?", announcementID).Delete(&entity.TicketAnnouncementTag{}).Error; err != nil {
			return err
		}
		// 2. Add the clean collection if elements remain
		if len(tags) > 0 {
			if err := tx.Create(&tags).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *announcementRepository) Delete(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("announcement_id = ?", id).Delete(&entity.TicketAnnouncementTag{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND apartment_id = ?", id, apartmentID).Delete(&entity.Announcement{}).Error
	})
}
