package dto

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	FirstName string             `json:"first_name" binding:"required"`
	LastName  string             `json:"last_name" binding:"required"`
	Email     string             `json:"email" binding:"required,email"`
	Username  string             `json:"username" binding:"required"`
	Gender    *entity.GenderType `json:"gender" binding:"omitempty"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type PushUserToUnitRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
