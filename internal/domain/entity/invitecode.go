package entity

import (
	"time"

	"github.com/google/uuid"
)

type InviteCode struct {
	BaseModel

	ApartmentID uuid.UUID

	Code string

	ExpiresAt time.Time

	Apartment Apartment
}
