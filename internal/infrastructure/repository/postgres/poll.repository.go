package postgres

import (
	"apartment-manager-backend/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pollRepository struct {
	db *gorm.DB
}

func NewPollRepository(db *gorm.DB) *pollRepository {
	return &pollRepository{db: db}
}

func (r *pollRepository) Create(poll *entity.Poll) error {
	return r.db.Create(poll).Error
}

func (r *pollRepository) GetByID(id uuid.UUID) (*entity.Poll, error) {
	var poll entity.Poll
	err := r.db.Preload("Options").Where("id = ? AND deleted_at IS NULL", id).First(&poll).Error
	if err != nil {
		return nil, err
	}
	return &poll, nil
}

func (r *pollRepository) ListByApartmentID(apartmentID uuid.UUID) ([]entity.Poll, error) {
	var polls []entity.Poll
	err := r.db.Preload("Options").Where("apartment_id = ? AND deleted_at IS NULL", apartmentID).Order("created_at DESC").Find(&polls).Error
	return polls, err
}

func (r *pollRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&entity.Poll{}).Error
}

func (r *pollRepository) GetVote(userID, pollID uuid.UUID) (*entity.Vote, error) {
	var vote entity.Vote
	err := r.db.Joins("Join poll_options On poll_options.id = votes.option_id").
		Where("votes.user_id = ? AND poll_options.poll_id = ? AND votes.deleted_at IS NULL", userID, pollID).
		First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *pollRepository) CreateVote(vote *entity.Vote) error {
	return r.db.Create(vote).Error
}

func (r *pollRepository) UpdateVote(vote *entity.Vote) error {
	return r.db.Save(vote).Error
}

func (r *pollRepository) GetVotesCount(pollID uuid.UUID) (map[uuid.UUID]int64, error) {
	type Result struct {
		OptionID uuid.UUID
		Count    int64
	}
	var results []Result

	err := r.db.Model(&entity.Vote{}).
		Select("votes.option_id, count(votes.id) as count").
		Joins("Join poll_options On poll_options.id = votes.option_id").
		Where("poll_options.poll_id = ? AND votes.deleted_at IS NULL", pollID).
		Group("votes.option_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	votesMap := make(map[uuid.UUID]int64)
	for _, res := range results {
		votesMap[res.OptionID] = res.Count
	}
	return votesMap, nil
}
