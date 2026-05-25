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
	Create(ctx context.Context, ticket *dto.CreateTicketRequest) (*dto.CreateTicketResponse, error)

	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*dto.TicketBaseResponse, error)
	getByIDSuperAccess(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)

	GetByIDWithAllRelations(ctx context.Context, id uuid.UUID) (*dto.TicketBaseResponseWithAllRelations, error)

	List(ctx context.Context, filter dto.TicketFilterRequest) ([]dto.TicketBaseResponseWithCommentCount, error)

	Update(ctx context.Context, id uuid.UUID, req dto.UpdateTicketRequest) (*dto.TicketBaseResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TicketStatus) (*dto.TicketBaseResponse, error)
	GetUserTickets(ctx context.Context, userID string) ([]dto.TicketResponse, error)
}

type ticketService struct {
	repo domainRepo.TicketInterface
}

func NewTicketService(repo domainRepo.TicketInterface) TicketService {
	return &ticketService{repo: repo}
}

func (s *ticketService) Create(ctx context.Context, req *dto.CreateTicketRequest) (*dto.CreateTicketResponse, error) {
	rawBaseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if rawBaseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	str_baseUserID := rawBaseUserID.(string)
	baseUserID, err := uuid.Parse(str_baseUserID)

	if err != nil {
		return nil, service_error.ErrCommonParseStrToUUID
	}

	// if req.UserID == nil || baseUserID != *req.UserID {
	// 	return nil, service_error.ErrTicketUnauthorizedAccess
	// }
	ticket := &entity.Ticket{
		UserID:        &baseUserID,
		Title:         req.Title,
		Description:   req.Description,
		Body:          req.Body,
		Category:      req.Category,
		Accessibility: req.Accessibility,
		Status:        entity.TicketOpen,
	}
	err = s.repo.Create(ctx, ticket)
	return &dto.CreateTicketResponse{
		ID:          ticket.ID,
		UserID:      ticket.UserID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Body:        ticket.Body,
		Category:    string(ticket.Category),
		Status:      string(ticket.Status),
	}, err
}

func (s *ticketService) Delete(ctx context.Context, id uuid.UUID) error {
	rawBaseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if rawBaseUserID == nil {
		return service_error.ErrUserIDNotFoundInContext
	}
	str_baseUserID := rawBaseUserID.(string)
	baseUserID, err := uuid.Parse(str_baseUserID)

	if err != nil {
		return service_error.ErrCommonParseStrToUUID
	}

	rawRole := ctx.Value("role") // IT MUST BE EXIST!
	if rawRole == nil {
		return service_error.ErrUserRoleNotFoundInContext
	}
	str_role := rawRole.(string)
	role := entity.UserRole(str_role)

	ticket, err := s.getByIDSuperAccess(ctx, id)
	if err != nil {
		return err
	}
	if ticket.UserID == nil || baseUserID != *(ticket.UserID) {
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

func (s *ticketService) GetByID(ctx context.Context, id uuid.UUID) (*dto.TicketBaseResponse, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	rawBaseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if rawBaseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	str_baseUserID := rawBaseUserID.(string)
	baseUserID, err := uuid.Parse(str_baseUserID)

	if err != nil {
		return nil, service_error.ErrCommonParseStrToUUID
	}

	rawRole := ctx.Value("role") // IT MUST BE EXIST!
	if rawRole == nil {
		return nil, service_error.ErrUserRoleNotFoundInContext
	}
	str_role := rawRole.(string)
	role := entity.UserRole(str_role)
	if ticket.Accessibility == entity.PrivateTicket && (ticket.UserID == nil || *ticket.UserID != baseUserID) {
		if !(role == entity.RoleManager || role == entity.RoleAdmin) {
			return nil, service_error.ErrTicketIsPrivate
		}

	}
	return dto.MapTicketToBaseResponse(ticket), nil
}

func (s *ticketService) getByIDSuperAccess(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	return ticket, nil
}

func (s *ticketService) GetByIDWithAllRelations(ctx context.Context, id uuid.UUID) (*dto.TicketBaseResponseWithAllRelations, error) {
	ticket, err := s.repo.GetByIDWithAllRelations(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	rawBaseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if rawBaseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	str_baseUserID := rawBaseUserID.(string)
	baseUserID, err := uuid.Parse(str_baseUserID)

	if err != nil {
		return nil, service_error.ErrCommonParseStrToUUID
	}

	rawRole := ctx.Value("role") // IT MUST BE EXIST!
	if rawRole == nil {
		return nil, service_error.ErrUserRoleNotFoundInContext
	}
	str_role := rawRole.(string)
	role := entity.UserRole(str_role)
	if ticket.Accessibility == entity.PrivateTicket && (ticket.UserID == nil || *ticket.UserID != baseUserID) {
		if !(role == entity.RoleManager || role == entity.RoleAdmin) {
			return nil, service_error.ErrTicketIsPrivate
		}
	}

	return &dto.TicketBaseResponseWithAllRelations{
		TicketBaseResponse: *dto.MapTicketToBaseResponse(ticket),
		Comments:           dto.MapCommentsToSliceResponse(ticket.Comments),
		User:               *dto.MapUserToUserResponse(&ticket.User),
	}, nil
}

// TODO : This list it's not fully
func (s *ticketService) List(ctx context.Context, filter dto.TicketFilterRequest) ([]dto.TicketBaseResponseWithCommentCount, error) {

	rawBaseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if rawBaseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	str_baseUserID := rawBaseUserID.(string)
	baseUserID, err := uuid.Parse(str_baseUserID)

	if err != nil {
		return nil, service_error.ErrCommonParseStrToUUID
	}

	rawRole := ctx.Value("role") // IT MUST BE EXIST!
	if rawRole == nil {
		return nil, service_error.ErrUserRoleNotFoundInContext
	}
	str_role := rawRole.(string)
	role := entity.UserRole(str_role)

	new_filter := domainRepo.TicketFilter{
		UserID:   filter.UserID,
		Status:   filter.Status,
		Category: filter.Category,
		Limit:    filter.Limit,
		Offset:   filter.Limit * (filter.Page - 1),
	}

	fl, err := s.repo.List(ctx, new_filter, baseUserID, role)
	if err != nil {
		return nil, err
	}

	return dto.MapTicketsToSliceResponseWithCount(fl), err
}

func (s *ticketService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTicketRequest) (*dto.TicketBaseResponse, error) {
	rawBaseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if rawBaseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	str_baseUserID := rawBaseUserID.(string)
	baseUserID, err := uuid.Parse(str_baseUserID)

	if err != nil {
		return nil, service_error.ErrCommonParseStrToUUID
	}

	ticket, err := s.getByIDSuperAccess(ctx, id) // it must retrieve a record for checking after that
	if err != nil {
		return nil, err
	}

	if ticket.UserID == nil || baseUserID != *(ticket.UserID) {
		return nil, service_error.ErrTicketUnauthorizedAccess
	}

	if _, err = s.repo.UpdateCategory(ctx, id, req.Category); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrTicketNotFound
		}
		return nil, err
	}

	var t *entity.Ticket
	if t, err = s.repo.UpdateContent(ctx, id, req.Title, req.Description, req.Body); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrTicketNotFound
		}
		return nil, err
	}

	return dto.MapTicketToBaseResponse(t), nil
}

func (s *ticketService) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TicketStatus) (*dto.TicketBaseResponse, error) {
	t, err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrTicketNotFound
		}
		return nil, err
	}

	return dto.MapTicketToBaseResponse(t), nil
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
