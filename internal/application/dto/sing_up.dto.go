package dto

import "apartment-manager-backend/internal/domain/entity"

type RegisterOutput struct {
	User         *entity.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	Message      string       `json:"message"`
}

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required,min=3,max=20"`
	LastName  string `json:"last_name" binding:"required,min=3,max=20"`
	Username  string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,len=11,numeric,startswith=09"`
	Password  string `json:"password" binding:"required,min=6,max=30"`
	Gender    string `json:"gender"`
}
