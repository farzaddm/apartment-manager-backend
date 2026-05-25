package dto

import (
	"apartment-manager-backend/internal/domain/entity"
	"time"
)

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

func MapInviteCodeToResponse(ic *entity.InviteCode) *InviteCodeResponse {
	if ic == nil {
		return nil
	}
	return &InviteCodeResponse{
		ID:          ic.ID.String(),
		ApartmentID: ic.ApartmentID.String(),
		UnitID:      ic.UnitID.String(),
		Code:        ic.Code,
		ExpiresAt:   ic.ExpiresAt,
	}
}

func MapInviteCodesToSliceResponse(codes []entity.InviteCode) []InviteCodeResponse {
	if codes == nil {
		return []InviteCodeResponse{}
	}
	res := make([]InviteCodeResponse, len(codes))
	for i, c := range codes {
		mapped := MapInviteCodeToResponse(&c)
		if mapped != nil {
			res[i] = *mapped
		}
	}
	return res
}
