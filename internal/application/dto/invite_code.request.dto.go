package dto

import (
	"time"
)

type CreateInviteRequest struct {
	ApartmentID string     `json:"apartment_id" binding:"required,uuid"`
	UnitID      string     `json:"unit_id" binding:"required,uuid"`
	ExpiresAt   *time.Time `json:"expires_at" binding:"required"`
}

type ValidateInviteCodeRequest struct {
	Code string `json:"code"`
}
