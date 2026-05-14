package entity

import "github.com/google/uuid"

type TicketAnnouncementTag struct {
	BaseModel

	TagID uuid.UUID

	TicketID       *uuid.UUID
	AnnouncementID *uuid.UUID

	Tag Tag

	Ticket       *Ticket
	Announcement *Announcement
}
