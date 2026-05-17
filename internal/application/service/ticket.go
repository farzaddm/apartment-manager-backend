package service

import (
	"apartment-manager-backend/internal/application/dto"
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"errors"

	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketService interface {
	Create(ctx context.Context, ticket *entity.Ticket) error

	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)

	GetByIDWithAllRelations(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)

	List(ctx context.Context, filter domainRepo.TicketFilter) ([]domainRepo.TicketWithCommentCount, error)

	Update(ctx context.Context, id uuid.UUID, req dto.UpdateTicketRequest) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TicketStatus) error
}

type ticketService struct {
	repo domainRepo.TicketInterface
}

func NewTicketService(repo domainRepo.TicketInterface) TicketService {
	return &ticketService{repo: repo}
}

func (s *ticketService) Create(ctx context.Context, ticket *entity.Ticket) error {
	return s.repo.Create(ctx, ticket)
}

func (s *ticketService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrTicketNotFound
		}
		return err
	}

	return nil
}

func (s *ticketService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	return ticket, nil
}

func (s *ticketService) GetByIDWithAllRelations(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
	ticket, err := s.repo.GetByIDWithAllRelations(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	return ticket, nil
}

func (s *ticketService) List(ctx context.Context, filter domainRepo.TicketFilter) ([]domainRepo.TicketWithCommentCount, error) {
	return s.repo.List(ctx, filter)
}

func (s *ticketService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTicketRequest) error {
	if err := s.repo.UpdateCategory(ctx, id, req.Category); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrTicketNotFound
		}
		return err
	}

	if err := s.repo.UpdateContent(ctx, id, req.Title, req.Description, req.Body); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrTicketNotFound
		}
		return err
	}

	return nil
}

func (s *ticketService) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TicketStatus) error {
	err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrTicketNotFound
		}
		return err
	}

	return nil
}
