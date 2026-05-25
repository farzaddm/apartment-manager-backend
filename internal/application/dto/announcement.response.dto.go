package dto

import (
	"apartment-manager-backend/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

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

func MapAnnouncementToResponse(a *entity.Announcement) *AnnouncementResponse {
	if a == nil {
		return nil
	}

	var tags = make([]TagResponse, 0)

	for i := range a.Tags {
		tags = append(tags, *MapTagToResponse(&a.Tags[i].Tag))
	}

	return &AnnouncementResponse{
		ID:          a.ID,
		ApartmentID: a.ApartmentID,
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
		Order:       a.Order,
		IsPinned:    a.IsPinned,
		ExpiredDate: a.ExpiredDate,
		Tags:        tags, //TODO : Fix Ann-Tag Relations
		CreatedAt:   a.CreatedAt,
	}
}

func MapAnnouncementsToResponseSlice(list []entity.Announcement) []AnnouncementResponse {
	if len(list) == 0 {
		return []AnnouncementResponse{}
	}

	result := make([]AnnouncementResponse, 0, len(list))

	for _, item := range list {
		result = append(result, *MapAnnouncementToResponse(&item))
	}

	return result
}
