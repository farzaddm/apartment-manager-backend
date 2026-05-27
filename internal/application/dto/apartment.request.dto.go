package dto

import (
	"apartment-manager-backend/internal/domain/entity"
)

type CreateApartmentRequest struct {
	Name       string                           `json:"name" binding:"required"`
	Province   string                           `json:"province" binding:"required"`
	City       string                           `json:"city" binding:"required"`
	Address    string                           `json:"address" binding:"required"`
	PostalCode string                           `json:"postal_code" binding:"required"`
	Units      []CreateUnitFromApartmentRequest `json:"units" binding:"required"`
}

type CreateUnitFromApartmentRequest struct {
	Floor      int    `json:"floor" binding:"required"`
	UnitNumber string `json:"unit_number" binding:"required"`
}

type UpdateApartmentRequest struct {
	Name       string `json:"name"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
}

func MapCreateUnitsToEntities(reqUnits []CreateUnitFromApartmentRequest) []entity.Unit {
	units := make([]entity.Unit, len(reqUnits))

	for i, u := range reqUnits {
		units[i] = entity.Unit{
			Floor:      u.Floor,
			UnitNumber: u.UnitNumber,
		}
	}

	return units
}
