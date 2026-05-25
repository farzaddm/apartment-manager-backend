package postgres

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type PollRepository interface {
	Create(poll *entity.Poll) error
	GetByID(id uuid.UUID) (*entity.Poll, error)
	ListByApartmentID(apartmentID uuid.UUID) ([]entity.Poll, error)
	Delete(id uuid.UUID) error

	GetVote(userID, pollID uuid.UUID) (*entity.Vote, error)
	CreateVote(vote *entity.Vote) error
	UpdateVote(vote *entity.Vote) error
	GetVotesCount(pollID uuid.UUID) (map[uuid.UUID]int64, error)
}
