package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type AnnouncementInterface interface {
	Create(ctx context.Context, announcement *entity.Announcement) error
	FindByIDAndApartment(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) (*entity.Announcement, error)
	Update(ctx context.Context, announcement *entity.Announcement) error
	ReplaceTags(ctx context.Context, announcementID uuid.UUID, tags []entity.TicketAnnouncementTag) error
	Delete(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) error
}
