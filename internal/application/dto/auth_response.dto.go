package dto

import "apartment-manager-backend/internal/domain/entity"

type RegisterOutput struct {
	User         *entity.User `json:"user,omitempty"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	Message      string       `json:"message"`
}

type LoginOutput = RegisterOutput

type VerifyOTPOutput struct {
	User         *entity.User `json:"user,omitempty"`
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	Message      string       `json:"message"`
}
