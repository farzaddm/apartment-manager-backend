package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserResponseDTO struct {
	ID              uuid.UUID  `json:"id"`
	CreatedAt       time.Time  `json:"created_at"`
	ApartmentID     *uuid.UUID `json:"apartment_id"`
	UnitID          *uuid.UUID `json:"unit_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Role            string     `json:"role"`
	Gender          string     `json:"gender"`
	ProfileImageURL *string    `json:"profile_image_url"`
}

type RegisterOutput struct {
	User         *UserResponseDTO `json:"user,omitempty"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	Message      string           `json:"message"`
}

type LoginOutput = RegisterOutput

type VerifyOTPOutput struct {
	User         *UserResponseDTO `json:"user,omitempty"`
	AccessToken  string           `json:"access_token,omitempty"`
	RefreshToken string           `json:"refresh_token,omitempty"`
	Message      string           `json:"message"`
}
