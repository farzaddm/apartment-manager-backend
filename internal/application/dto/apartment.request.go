package dto

import "github.com/google/uuid"

type CreateApartmentRequest struct {
	Name       string `json:"name" binding:"required"`
	Province   string `json:"province" binding:"required"`
	City       string `json:"city" binding:"required"`
	Address    string `json:"address" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
}

type CreateApartmentResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
}

type UpdateApartmentRequest struct {
	Name       string `json:"name"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
}
