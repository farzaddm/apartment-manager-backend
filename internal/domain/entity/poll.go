package entity

import (
	"time"

	"github.com/google/uuid"
)

type Poll struct {
	BaseModel

	ApartmentID   uuid.UUID
	Title         string
	Description   string
	ExpiresAt     *time.Time
	IsVotesPublic bool

	Apartment Apartment
}
