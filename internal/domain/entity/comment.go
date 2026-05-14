package entity

import "github.com/google/uuid"

type Comment struct {
	BaseModel

	UserID   *uuid.UUID
	TicketID uuid.UUID

	Body string

	User   User
	Ticket Ticket
}
