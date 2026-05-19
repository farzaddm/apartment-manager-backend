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
	Create(ctx context.Context, baseUserID uuid.UUID, ticket *entity.Ticket) error

	Delete(ctx context.Context, baseUserID uuid.UUID, id uuid.UUID, role entity.UserRole) error

	GetByID(ctx context.Context, baseUserID uuid.UUID, id uuid.UUID, role entity.UserRole) (*entity.Ticket, error)

	GetByIDWithAllRelations(ctx context.Context, baseUserID uuid.UUID, id uuid.UUID, role entity.UserRole) (*entity.Ticket, error)

	List(ctx context.Context, baseUserID uuid.UUID, filter dto.TicketFilterRequest, role entity.UserRole) ([]domainRepo.TicketWithCommentCount, error)

	Update(ctx context.Context, baseUserID uuid.UUID, id uuid.UUID, req dto.UpdateTicketRequest) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TicketStatus) error
	GetUserTickets(ctx context.Context, userID string) ([]dto.TicketResponse, error)
}

type ticketService struct {
	repo domainRepo.TicketInterface
}

func NewTicketService(repo domainRepo.TicketInterface) TicketService {
	return &ticketService{repo: repo}
}

func (s *ticketService) Create(ctx context.Context, baseUserID uuid.UUID, ticket *entity.Ticket) error {
	//TODO: nil situation
	if ticket.UserID != nil && baseUserID == *ticket.UserID {
		return service_error.ErrTicketUnauthorizedAccess
	}
	return s.repo.Create(ctx, ticket)
}

func (s *ticketService) Delete(ctx context.Context, baseUserID, id uuid.UUID, role entity.UserRole) error {
	ticket, err := s.GetByID(ctx, baseUserID, id, role)
	if err != nil {
		return err
	}
	//TODO: nil situation
	if ticket.UserID != nil && baseUserID == *ticket.UserID {
		if !(role == entity.RoleManager || role == entity.RoleAdmin) {
			return service_error.ErrTicketUnauthorizedAccess
		}
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrTicketNotFound
		}
		return err
	}

	return nil
}

func (s *ticketService) GetByID(ctx context.Context, baseUserID uuid.UUID, id uuid.UUID, role entity.UserRole) (*entity.Ticket, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}
	if ticket.Accessability == entity.PrivateTicket && (ticket.UserID == nil || *ticket.UserID != baseUserID) {
		if !(role == entity.RoleManager || role == entity.RoleAdmin) {
			return nil, service_error.ErrTicketIsPrivate
		}

	}
	return ticket, nil
}

func (s *ticketService) GetByIDWithAllRelations(ctx context.Context, id uuid.UUID, baseUserID uuid.UUID, role entity.UserRole) (*entity.Ticket, error) {
	ticket, err := s.repo.GetByIDWithAllRelations(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}
	if ticket.Accessability == entity.PrivateTicket && (ticket.UserID == nil || *ticket.UserID != baseUserID) {
		if !(role == entity.RoleManager || role == entity.RoleAdmin) {
			return nil, service_error.ErrTicketIsPrivate
		}
	}

	return ticket, nil
}

func (s *ticketService) List(ctx context.Context, baseUserID uuid.UUID, filter dto.TicketFilterRequest, role entity.UserRole) ([]domainRepo.TicketWithCommentCount, error) {
	new_filter := domainRepo.TicketFilter{
		UserID:   filter.UserID,
		Status:   filter.Status,
		Category: filter.Category,
		Limit:    filter.Limit,
		Offset:   filter.Limit * (filter.Page - 1),
	}
	fl, err := s.repo.List(ctx, new_filter)
	if err != nil {
		return nil, err
	}
	new_fl := make([]domainRepo.TicketWithCommentCount, 0)
	for i := range fl {
		if fl[i].Accessability == entity.PublicTicket ||
			(fl[i].UserID != nil && *fl[i].UserID == baseUserID) ||
			role == entity.RoleManager || role == entity.RoleAdmin {
			new_fl = append(new_fl, fl[i])

		}
	}
	return new_fl, err
}

func (s *ticketService) Update(ctx context.Context, baseUserID, id uuid.UUID, req dto.UpdateTicketRequest) error {
	ticket, err := s.GetByID(ctx, baseUserID, id, entity.RoleAdmin) // it must retrieve a record for checking after that
	if err != nil {
		return err
	}
	//TODO: nil situation
	if ticket.UserID != nil && baseUserID == *ticket.UserID {
		return service_error.ErrTicketUnauthorizedAccess
	}

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

func MapTicketsToResponse(domainTickets []entity.Ticket) []dto.TicketResponse {
	responses := make([]dto.TicketResponse, len(domainTickets))

	for i, t := range domainTickets {
		var tagIDs []string
		for _, ticketTag := range t.Tags {
			tagIDs = append(tagIDs, ticketTag.TagID.String())
		}

		var uID uuid.UUID
		if t.UserID != nil {
			uID = *t.UserID
		}

		responses[i] = dto.TicketResponse{
			ID:          t.ID, // Inherited from your BaseModel
			UserID:      uID,
			Title:       t.Title,
			Description: t.Description,
			Body:        t.Body,
			Category:    string(t.Category),
			Status:      string(t.Status),
			Tags:        tagIDs,
		}
	}

	return responses
}

func (s *ticketService) GetUserTickets(ctx context.Context, userID string) ([]dto.TicketResponse, error) {
	rawTickets, err := s.repo.GetTicketsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	tickets := MapTicketsToResponse(rawTickets)
	return tickets, nil
}
