package dto

import (
	"apartment-manager-backend/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type CreateAnnouncementRequest struct {
	Title       string                   `json:"title" binding:"required,max=255"`
	Description string                   `json:"description"`
	Body        string                   `json:"body" binding:"required"`
	Order       entity.AnnouncementOrder `json:"order"`
	IsPinned    bool                     `json:"is_pinned"`
	ExpiredDate *time.Time               `json:"expired_date"`
	TagIDs      []uuid.UUID              `json:"tag_ids"`
}

type UpdateAnnouncementRequest struct {
	Title       string                   `json:"title" binding:"required,max=255"`
	Description string                   `json:"description"`
	Body        string                   `json:"body" binding:"required"`
	Order       entity.AnnouncementOrder `json:"order"`
	IsPinned    bool                     `json:"is_pinned"`
	ExpiredDate *time.Time               `json:"expired_date"`
	TagIDs      []uuid.UUID              `json:"tag_ids"`
}

type AnnouncementResponse struct {
	ID          uuid.UUID                `json:"id"`
	ApartmentID uuid.UUID                `json:"apartment_id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Body        string                   `json:"body"`
	Order       entity.AnnouncementOrder `json:"order"`
	IsPinned    bool                     `json:"is_pinned"`
	ExpiredDate *time.Time               `json:"expired_date,omitempty"`
	Tags        []TagResponse            `json:"tags"`
	CreatedAt   time.Time                `json:"created_at"`
}
