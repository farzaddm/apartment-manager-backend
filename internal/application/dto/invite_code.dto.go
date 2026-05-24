package dto

import "time"

type CreateInviteRequest struct {
	ApartmentID string     `json:"apartment_id" binding:"required,uuid"`
	UnitID      string     `json:"unit_id" binding:"required,uuid"`
	ExpiresAt   *time.Time `json:"expires_at" binding:"required"`
}

type InviteCodeResponse struct {
	ID          string    `json:"id"`
	ApartmentID string    `json:"apartment_id"`
	UnitID      string    `json:"unit_id"`
	Code        string    `json:"code"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type ValidateInviteCodeRequest struct {
	Code string
}

type ValidateInviteCodeResponse struct {
	ApartmentID string `json:"apartment_id"`
	UnitID      string `json:"unit_id"`
}
