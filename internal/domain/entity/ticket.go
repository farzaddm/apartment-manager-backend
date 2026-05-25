package entity

import "github.com/google/uuid"

type Ticket struct {
	BaseModel

	UserID *uuid.UUID

	Title       string
	Description string
	Body        string

	Category TicketCategory

	Status TicketStatus

	Accessibility TicketAccessibility

	User User

	Comments []Comment

	Tags []TicketAnnouncementTag
}
