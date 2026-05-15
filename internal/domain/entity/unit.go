package entity

import "github.com/google/uuid"

type Unit struct {
	BaseModel

	ApartmentID uuid.UUID
	UserID      *uuid.UUID

	UnitNumber string
	Floor      int

	Apartment Apartment
	User      User
}
