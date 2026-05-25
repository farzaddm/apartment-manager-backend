package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePollRequest struct {
	Title         string     `json:"title" binding:"required"`
	Description   string     `json:"description"`
	ExpiresAt     *time.Time `json:"expires_at"`
	IsVotesPublic bool       `json:"is_votes_public"`
	Options       []string   `json:"options" binding:"required,min=2"`
}

type VoteRequest struct {
	OptionID uuid.UUID `json:"option_id" binding:"required"`
}

type PollOptionResponse struct {
	ID         uuid.UUID `json:"id"`
	Text       string    `json:"text"`
	VotesCount int64     `json:"votes_count"`
}

type PollResponse struct {
	ID                uuid.UUID            `json:"id"`
	Title             string               `json:"title"`
	Description       string               `json:"description"`
	ExpiresAt         *time.Time           `json:"expires_at"`
	IsVotesPublic     bool                 `json:"is_votes_public"`
	TotalVotes        int64                `json:"total_votes"`
	Options           []PollOptionResponse `json:"options"`
	UserVotedOptionID *uuid.UUID           `json:"user_voted_option_id,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
}
