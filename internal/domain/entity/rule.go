package entity

import (
	"github.com/google/uuid"
)

type Rule struct {
	BaseModel

	ApartmentID uuid.UUID    `json:"apartment_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Category    RuleCategory `json:"category"`

	Apartment Apartment `json:"apartment,omitempty"`
}
