package dto

import "github.com/google/uuid"

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
