package dto

import (
	"apartment-manager-backend/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type UnitResponse struct {
	ID          uuid.UUID  `json:"id"`
	ApartmentID uuid.UUID  `json:"apartment_id"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	UnitNumber  string     `json:"unit_number"`
	Floor       int        `json:"floor"`
	CreatedAt   time.Time  `json:"created_at"`

	Apartment *ApartmentResponse `json:"apartment,omitempty"`
	User      *UserResponse      `json:"user,omitempty"`
}

func MapUnitToResponse(unit *entity.Unit) *UnitResponse {
	if unit == nil {
		return nil
	}

	res := &UnitResponse{
		ID:          unit.ID,
		ApartmentID: unit.ApartmentID,
		UserID:      unit.UserID,
		UnitNumber:  unit.UnitNumber,
		Floor:       unit.Floor,
		CreatedAt:   unit.CreatedAt,
	}

	if unit.Apartment.ID != uuid.Nil {
		res.Apartment = &ApartmentResponse{
			ID:   unit.Apartment.ID,
			Name: unit.Apartment.Name,
		}
	}

	if unit.UserID != nil && unit.User.ID != (uuid.UUID{}) {
		res.User = MapUserToUserResponse(&unit.User)
	}

	return res
}

func MapUnitToSliceResponse(units []entity.Unit) []UnitResponse {
	responses := make([]UnitResponse, len(units))
	for i, unit := range units {
		responses[i] = *MapUnitToResponse(&unit)
	}
	return responses
}
