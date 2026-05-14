package entity

import "github.com/google/uuid"

type Rule struct {
	BaseModel

	ApartmentID uuid.UUID

	Category string

	Apartment Apartment

	Items []RuleItem
}
