package entity

import (
	"github.com/google/uuid"
)

type PollOption struct {
	BaseModel

	PollID uuid.UUID
	Text   string

	Poll Poll
}
