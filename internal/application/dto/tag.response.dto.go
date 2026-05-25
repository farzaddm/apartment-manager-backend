package dto

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func MapTagToResponse(tag *entity.Tag) *TagResponse {
	if tag == nil {
		return nil
	}

	return &TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func MapTagsToSliceResponse(tags []entity.Tag) []TagResponse {
	if len(tags) == 0 {
		return []TagResponse{}
	}

	res := make([]TagResponse, 0, len(tags))
	for i := range tags {
		res = append(res, *MapTagToResponse(&tags[i]))
	}

	return res
}
