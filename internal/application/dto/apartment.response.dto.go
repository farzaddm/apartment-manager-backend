package dto

import (
	"apartment-manager-backend/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type CreateApartmentResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	CreatedAt  time.Time `json:"created_at"`
}

type ApartmentResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	CreatedAt  time.Time `json:"created_at"`
}

type ApartmentResponseWithUsersOfUnits struct {
	ApartmentResponse
	Units []UnitResponse `json:"units"`
}

type ApartmentResponseWithUsers struct {
	ApartmentResponse
	Users []UserResponse `json:"users"`
}

type ApartmentResponseWithAnnouncements struct {
	ApartmentResponse
	Announcements []AnnouncementResponse `json:"announcements"`
}

type ApartmentResponseWithRules struct {
	ApartmentResponse
	Rules []RuleResponse `json:"rules"`
}

type ApartmentResponseWithInviteCodes struct {
	ApartmentResponse
	InviteCodes []InviteCodeResponse `json:"invite_codes"`
}

func MapApartmentToResponse(apartment *entity.Apartment) *ApartmentResponse {
	if apartment == nil {
		return nil
	}
	return &ApartmentResponse{
		ID:         apartment.ID,
		Name:       apartment.Name,
		Province:   apartment.Province,
		City:       apartment.City,
		Address:    apartment.Address,
		PostalCode: apartment.PostalCode,
		CreatedAt:  apartment.CreatedAt,
	}
}
