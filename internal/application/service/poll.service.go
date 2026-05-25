package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PollService interface {
	CreatePoll(apartmentID uuid.UUID, req dto.CreatePollRequest) (*dto.PollResponse, error)
	ListPolls(apartmentID uuid.UUID, userID uuid.UUID, userApartmentID *uuid.UUID, userRole string) ([]dto.PollResponse, error)
	GetPollDetails(pollID uuid.UUID, userID uuid.UUID, userApartmentID *uuid.UUID, userRole string) (*dto.PollResponse, error)
	DeletePoll(pollID uuid.UUID) error
	CastVote(pollID uuid.UUID, userID uuid.UUID, userApartmentID *uuid.UUID, userRole string, optionID uuid.UUID) error
	RevokeVote(pollID uuid.UUID, userID uuid.UUID, userApartmentID *uuid.UUID, userRole string) error
}

type pollService struct {
	pollRepo postgres.PollRepository
}

func NewPollService(pollRepo postgres.PollRepository) PollService {
	return &pollService{pollRepo: pollRepo}
}

func (s *pollService) CreatePoll(apartmentID uuid.UUID, req dto.CreatePollRequest) (*dto.PollResponse, error) {
	poll := &entity.Poll{
		ApartmentID:   apartmentID,
		Title:         req.Title,
		Description:   req.Description,
		ExpiresAt:     req.ExpiresAt,
		IsVotesPublic: req.IsVotesPublic,
	}

	for _, optText := range req.Options {
		poll.Options = append(poll.Options, entity.PollOption{
			Text: optText,
		})
	}

	if err := s.pollRepo.Create(poll); err != nil {
		return nil, err
	}

	var optionsResp []dto.PollOptionResponse
	for _, opt := range poll.Options {
		optionsResp = append(optionsResp, dto.PollOptionResponse{
			ID:         opt.ID,
			Text:       opt.Text,
			VotesCount: 0,
		})
	}

	return &dto.PollResponse{
		ID:            poll.ID,
		Title:         poll.Title,
		Description:   poll.Description,
		ExpiresAt:     poll.ExpiresAt,
		IsVotesPublic: poll.IsVotesPublic,
		TotalVotes:    0,
		Options:       optionsResp,
		CreatedAt:     poll.CreatedAt,
	}, nil
}

func (s *pollService) CastVote(pollID uuid.UUID, userID uuid.UUID, userApartmentID *uuid.UUID, userRole string, optionID uuid.UUID) error {
	poll, err := s.pollRepo.GetByID(pollID)
	if err != nil {
		return errors.New("poll not found")
	}

	if userRole != "admin" && (userApartmentID == nil || *userApartmentID != poll.ApartmentID) {
		return errors.New("access denied: you do not belong to this apartment")
	}

	if poll.ExpiresAt != nil && time.Now().After(*poll.ExpiresAt) {
		return errors.New("this poll has expired")
	}

	validOption := false
	for _, opt := range poll.Options {
		if opt.ID == optionID {
			validOption = true
			break
		}
	}
	if !validOption {
		return errors.New("selected option does not belong to this poll")
	}

	existingVote, _ := s.pollRepo.GetVote(userID, pollID)
	if existingVote != nil {
		if existingVote.OptionID == optionID {
			return nil
		}
		existingVote.OptionID = optionID
		return s.pollRepo.UpdateVote(existingVote)
	}

	newVote := &entity.Vote{
		UserID:   userID,
		OptionID: optionID,
	}
	return s.pollRepo.CreateVote(newVote)
}

func (s *pollService) RevokeVote(pollID uuid.UUID, userID uuid.UUID, userApartmentID *uuid.UUID, userRole string) error {
	poll, err := s.pollRepo.GetByID(pollID)
	if err != nil {
		return errors.New("poll not found")
	}

	if userRole != "admin" && (userApartmentID == nil || *userApartmentID != poll.ApartmentID) {
		return errors.New("access denied: you do not belong to this apartment")
	}

	if poll.ExpiresAt != nil && time.Now().After(*poll.ExpiresAt) {
		return errors.New("poll has expired and votes cannot be revoked")
	}

	existingVote, err := s.pollRepo.GetVote(userID, pollID)
	if err != nil || existingVote == nil {
		return errors.New("no vote found for this poll")
	}

	existingVote.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	return s.pollRepo.UpdateVote(existingVote)
}

func (s *pollService) ListPolls(apartmentID uuid.UUID, userID uuid.UUID, userApartmentID *uuid.UUID, userRole string) ([]dto.PollResponse, error) {
	if userRole != "admin" && (userApartmentID == nil || *userApartmentID != apartmentID) {
		return nil, errors.New("access denied: you do not belong to this apartment")
	}

	polls, err := s.pollRepo.ListByApartmentID(apartmentID)
	if err != nil {
		return nil, err
	}

	var response []dto.PollResponse

	for _, poll := range polls {
		votesCount, _ := s.pollRepo.GetVotesCount(poll.ID)
		userVote, _ := s.pollRepo.GetVote(userID, poll.ID)

		var userVotedOptionID *uuid.UUID = nil
		if userVote != nil {
			idCopy := userVote.OptionID
			userVotedOptionID = &idCopy
		}

		isAuthorizedToSeeResults := poll.IsVotesPublic || userRole == "manager" || userRole == "admin"

		var optionsResp []dto.PollOptionResponse
		var totalVotes int64 = 0

		for _, opt := range poll.Options {
			count := votesCount[opt.ID]
			totalVotes += count

			if !isAuthorizedToSeeResults {
				count = -1
			}

			optionsResp = append(optionsResp, dto.PollOptionResponse{
				ID:         opt.ID,
				Text:       opt.Text,
				VotesCount: count,
			})
		}

		finalTotalVotes := totalVotes
		if !isAuthorizedToSeeResults {
			finalTotalVotes = -1
		}

		response = append(response, dto.PollResponse{
			ID:                poll.ID,
			Title:             poll.Title,
			Description:       poll.Description,
			ExpiresAt:         poll.ExpiresAt,
			IsVotesPublic:     poll.IsVotesPublic,
			TotalVotes:        finalTotalVotes,
			Options:           optionsResp,
			UserVotedOptionID: userVotedOptionID,
			CreatedAt:         poll.CreatedAt,
		})
	}

	return response, nil
}

func (s *pollService) GetPollDetails(pollID uuid.UUID, userID uuid.UUID, userApartmentID *uuid.UUID, userRole string) (*dto.PollResponse, error) {
	poll, err := s.pollRepo.GetByID(pollID)
	if err != nil {
		return nil, errors.New("poll not found")
	}

	if userRole != "admin" && (userApartmentID == nil || *userApartmentID != poll.ApartmentID) {
		return nil, errors.New("access denied: you do not belong to this apartment")
	}

	isAuthorizedToSeeResults := poll.IsVotesPublic || userRole == "manager" || userRole == "admin"

	votesCount, _ := s.pollRepo.GetVotesCount(poll.ID)
	userVote, _ := s.pollRepo.GetVote(userID, poll.ID)

	var userVotedOptionID *uuid.UUID = nil
	if userVote != nil {
		idCopy := userVote.OptionID
		userVotedOptionID = &idCopy
	}

	var optionsResp []dto.PollOptionResponse
	var totalVotes int64 = 0

	for _, opt := range poll.Options {
		count := votesCount[opt.ID]
		totalVotes += count

		if !isAuthorizedToSeeResults {
			count = -1
		}

		optionsResp = append(optionsResp, dto.PollOptionResponse{
			ID:         opt.ID,
			Text:       opt.Text,
			VotesCount: count,
		})
	}

	finalTotalVotes := totalVotes
	if !isAuthorizedToSeeResults {
		finalTotalVotes = -1
	}

	return &dto.PollResponse{
		ID:                poll.ID,
		Title:             poll.Title,
		Description:       poll.Description,
		ExpiresAt:         poll.ExpiresAt,
		IsVotesPublic:     poll.IsVotesPublic,
		TotalVotes:        finalTotalVotes,
		Options:           optionsResp,
		UserVotedOptionID: userVotedOptionID,
		CreatedAt:         poll.CreatedAt,
	}, nil
}

func (s *pollService) DeletePoll(pollID uuid.UUID) error {
	return s.pollRepo.Delete(pollID)
}
