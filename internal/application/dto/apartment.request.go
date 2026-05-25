package dto

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type CreateApartmentRequest struct {
	Name       string `json:"name" binding:"required"`
	Province   string `json:"province" binding:"required"`
	City       string `json:"city" binding:"required"`
	Address    string `json:"address" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
}

type UpdateApartmentRequest struct {
	Name       string `json:"name"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
}

type ApartmentResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Province     string    `json:"province"`
	City         string    `json:"city"`
	Address      string    `json:"address"`
	PostalCode   string    `json:"postal_code"`
	CreationDate string    `json:"creation_date"`
}

type ApartmentResponseWithUsers struct {
	ApartmentResponse
	Users []UserResponse
}

type ApartmentResponseWithAnnouncements struct {
	ApartmentResponse
	Announcements []AnnouncementResponse
}

type ApartmentResponseWithRules struct {
	ApartmentResponse
	Rules []RuleResponse
}
type ApartmentResponseWithInviteCodes struct {
	ApartmentResponse
	InviteCodes []InviteCodeResponse
}

func MapApartmentToResponse(apartment *entity.Apartment) *ApartmentResponse {
	if apartment == nil {
		return nil
	}
	return &ApartmentResponse{
		ID:           apartment.ID,
		Name:         apartment.Name,
		Province:     apartment.Province,
		City:         apartment.City,
		Address:      apartment.Address,
		PostalCode:   apartment.PostalCode,
		CreationDate: apartment.CreatedAt.Format("2006-01-02"),
	}
}
