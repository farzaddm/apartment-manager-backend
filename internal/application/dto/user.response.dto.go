package dto

import (
	"apartment-manager-backend/internal/domain/entity"
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
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Username  string             `json:"username"`
	Gender    *entity.GenderType `json:"gender"`
	Role      entity.UserRole    `json:"role"`
	Phone     string             `json:"phone"`
}

// TODO : Hi :)
type UserResponseWithAllRelations struct {
	UserResponse
	Apartment *ApartmentResponse   `json:"apartment"`
	Comments  []XXXComment         `json:"comment"`
	Tickets   []TicketBaseResponse `json:"ticket"`
}

func MapUserToUserResponse(u *entity.User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
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
	if users == nil {
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
