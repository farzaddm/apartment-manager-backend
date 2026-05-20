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
	Create(ctx context.Context, ticket *dto.CreateTicketRequest) error

	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)
	GetByIDSuperAccess(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)

	GetByIDWithAllRelations(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)

	List(ctx context.Context, filter dto.TicketFilterRequest) ([]domainRepo.TicketWithCommentCount, error)

	Update(ctx context.Context, id uuid.UUID, req dto.UpdateTicketRequest) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TicketStatus) error
	GetUserTickets(ctx context.Context, userID string) ([]dto.TicketResponse, error)
}

type ticketService struct {
	repo domainRepo.TicketInterface
}

func NewTicketService(repo domainRepo.TicketInterface) TicketService {
	return &ticketService{repo: repo}
}

func (s *ticketService) Create(ctx context.Context, req *dto.CreateTicketRequest) error {
	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return service_error.ErrUserIDNotFoundInContext
	}

	//TODO: nil situation
	if req.UserID != nil && baseUserID == *req.UserID {
		return service_error.ErrTicketUnauthorizedAccess
	}
	ticket := &entity.Ticket{
		UserID:        req.UserID,
		Title:         req.Title,
		Description:   req.Description,
		Body:          req.Body,
		Category:      req.Category,
		Accessability: req.Accessability,
	}
	return s.repo.Create(ctx, ticket)
}

func (s *ticketService) Delete(ctx context.Context, id uuid.UUID) error {
	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return service_error.ErrUserIDNotFoundInContext
	}
	role := ctx.Value("role") // IT MUST BE EXIST!
	if role == nil {
		return service_error.ErrUserRoleNotFoundInContext
	}

	//TODO: nil situation
	ticket, err := s.GetByIDSuperAccess(ctx, id)
	if err != nil {
		return err
	}
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

func (s *ticketService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	role := ctx.Value("role") // IT MUST BE EXIST!
	if role == nil {
		return nil, service_error.ErrUserRoleNotFoundInContext
	}
	if ticket.Accessability == entity.PrivateTicket && (ticket.UserID == nil || *ticket.UserID != baseUserID) {
		if !(role == entity.RoleManager || role == entity.RoleAdmin) {
			return nil, service_error.ErrTicketIsPrivate
		}

	}
	return ticket, nil
}

func (s *ticketService) GetByIDSuperAccess(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
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

	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	role := ctx.Value("role") // IT MUST BE EXIST!
	if role == nil {
		return nil, service_error.ErrUserRoleNotFoundInContext
	}
	if ticket.Accessability == entity.PrivateTicket && (ticket.UserID == nil || *ticket.UserID != baseUserID) {
		if !(role == entity.RoleManager || role == entity.RoleAdmin) {
			return nil, service_error.ErrTicketIsPrivate
		}
	}

	return ticket, nil
}

func (s *ticketService) List(ctx context.Context, filter dto.TicketFilterRequest) ([]domainRepo.TicketWithCommentCount, error) {
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

	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	role := ctx.Value("role") // IT MUST BE EXIST!
	if role == nil {
		return nil, service_error.ErrUserRoleNotFoundInContext
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

func (s *ticketService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTicketRequest) error {
	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return service_error.ErrUserIDNotFoundInContext
	}

	ticket, err := s.GetByIDSuperAccess(ctx, id) // it must retrieve a record for checking after that
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
