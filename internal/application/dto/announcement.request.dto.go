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
