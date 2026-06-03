package dto

import (
	"apartment-manager-backend/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type UserProfileResponse struct {
	ID              uuid.UUID  `json:"id"`
	CreatedAt       time.Time  `json:"created_at"`
	ApartmentID     *uuid.UUID `json:"apartment_id"`
	UnitNumber      string     `json:"unit_number"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Role            string     `json:"role"`
	Gender          string     `json:"gender"`
	ProfileImageURL *string    `json:"profile_image_url"`
}

type UpdateProfileResponse struct {
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Username  string             `json:"username"`
	Gender    *entity.GenderType `json:"gender,omitempty"`
}

type UploadImageResponse struct {
	ProfileImageURL *string `json:"profile_image_url"`
}

type UserResponse struct {
	ID              uuid.UUID          `json:"user_id"`
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
	Email           string             `json:"email"`
	Username        string             `json:"username"`
	Gender          *entity.GenderType `json:"gender"`
	Role            entity.UserRole    `json:"role"`
	Phone           string             `json:"phone"`
	Unit            *UnitResponse      `json:"unit,omitempty"`
	ProfileImageURL *string            `json:"profile_image_url"`
	CreatedAt       time.Time          `json:"created_at"`
}

type UserResponseWithAllRelations struct {
	UserResponse
	Apartment *ApartmentResponse   `json:"apartment"`
	Comments  []CommentResponse    `json:"comment"`
	Tickets   []TicketBaseResponse `json:"ticket"`
}

func MapUserToUserResponse(u *entity.User) *UserResponse {
	if u == nil {
		return nil
	}

	return &UserResponse{
		ID:              u.ID,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Email:           u.Email,
		Username:        u.Username,
		Gender:          u.Gender,
		Role:            u.Role,
		Phone:           u.Phone,
		Unit:            MapUnitToResponse(u.Unit),
		ProfileImageURL: u.ProfileImageURL,
		CreatedAt:       u.CreatedAt,
	}
}

func MapUsersToSliceResponse(users []entity.User) []UserResponse {
	if len(users) == 0 {
		return []UserResponse{}
	}
	res := make([]UserResponse, len(users))
	for i, u := range users {
		mapped := MapUserToUserResponse(&u)
		if mapped != nil {
			res[i] = *mapped
		}
	}
	return res
}
