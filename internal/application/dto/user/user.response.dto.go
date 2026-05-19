package userdto

import "apartment-manager-backend/internal/domain/entity"

type UserProfileResponse struct {
	UserId    string             `json:"user_id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	Gender    *entity.GenderType `json:"gender"`
	Role      entity.UserRole    `json:"role"`
}
