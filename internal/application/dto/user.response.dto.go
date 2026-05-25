package dto

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type UserProfileResponse struct {
	UserId    string             `json:"user_id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Username  string             `json:"username"`
	Gender    *entity.GenderType `json:"gender"`
	Role      entity.UserRole    `json:"role"`
}

type UserResponse struct {
	ID        uuid.UUID          `json:"user_id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Username  string             `json:"username"`
	Gender    *entity.GenderType `json:"gender"`
	Role      entity.UserRole    `json:"role"`
	Phone     string             `json:"phone"`
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
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Username:  u.Username,
		Gender:    u.Gender,
		Role:      u.Role,
		Phone:     u.Phone,
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
