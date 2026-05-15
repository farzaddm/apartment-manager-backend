package entity

import (
	"github.com/google/uuid"
)

type Vote struct {
	BaseModel

	UserID   uuid.UUID
	OptionID uuid.UUID

	User       User
	PollOption PollOption
}
