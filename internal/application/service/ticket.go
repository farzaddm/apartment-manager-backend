package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"

	"context"

	"github.com/google/uuid"
)

type TicketService interface {
	Create(ctx context.Context, ticket *entity.Ticket) error

	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)

	GetByIDWithAllRelations(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)

	List(ctx context.Context, filter domainRepo.TicketFilter) ([]domainRepo.TicketWithCommentCount, error)

	Update(ctx context.Context, id uuid.UUID, req dto.UpdateTicketRequest) error
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
	return s.repo.Delete(ctx, id)
}

func (s *ticketService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ticketService) GetByIDWithAllRelations(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
	return s.repo.GetByIDWithAllRelations(ctx, id)
}

func (s *ticketService) List(ctx context.Context, filter domainRepo.TicketFilter) ([]domainRepo.TicketWithCommentCount, error) {
	return s.repo.List(ctx, filter)
}

func (s *ticketService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTicketRequest) error {
	if err := s.repo.UpdateStatus(ctx, id, req.Status); err != nil {
		return err
	}

	if err := s.repo.UpdateCategory(ctx, id, req.Category); err != nil {
		return err
	}

	if err := s.repo.UpdateContent(ctx, id, req.Title, req.Description, req.Body); err != nil {
		return err
	}

	return nil
}
