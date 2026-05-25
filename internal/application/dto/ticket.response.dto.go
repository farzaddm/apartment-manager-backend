package dto

import "github.com/google/uuid"

type TicketResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
	Tags        []string  `json:"tags,omitempty"` // Returns an array of clean Tag IDs string
}
