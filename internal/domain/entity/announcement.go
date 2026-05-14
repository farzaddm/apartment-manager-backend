package entity

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	BaseModel

	ApartmentID uuid.UUID

	Title       string
	Description string
	Body        string

	ExpiredDate *time.Time
	Apartment   Apartment

	Tags []TicketAnnouncementTag
}
