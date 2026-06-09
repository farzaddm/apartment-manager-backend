package service

import (
	"apartment-manager-backend/internal/application/dto"
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/internal/domain/constant"
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"errors"

	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketService interface {
	Create(ctx context.Context, tokenKeys *dto.TokenKeys, ticket *dto.CreateTicketRequest) (*dto.CreateTicketResponse, error)

	Delete(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) error

	GetByID(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) (*dto.TicketBaseResponse, error)

	GetByIDWithAllRelations(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) (*dto.TicketBaseResponseWithAllRelations, error)

	List(ctx context.Context, tokenKeys *dto.TokenKeys, filter dto.TicketFilterRequest) ([]dto.TicketBaseResponseWithCommentCount, error)

	Update(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID, req dto.UpdateTicketRequest) (*dto.TicketBaseResponse, error)
	UpdateStatus(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID, status entity.TicketStatus) (*dto.TicketBaseResponse, error)
	GetUserTickets(ctx context.Context, userID string) ([]dto.TicketResponse, error)

	getByIDSuperAccess(ctx context.Context /*TODO : tokenKeys *dto.TokenKeys,*/, id uuid.UUID) (*entity.Ticket, error)
}

type ticketService struct {
	repo    domainRepo.TicketInterface
	tagRepo domainRepo.TagInterface
}

func NewTicketService(repo domainRepo.TicketInterface, tagRepo domainRepo.TagInterface) TicketService {
	return &ticketService{repo: repo, tagRepo: tagRepo}
}

func (s *ticketService) Create(ctx context.Context, tokenKeys *dto.TokenKeys, req *dto.CreateTicketRequest) (*dto.CreateTicketResponse, error) {

	// if req.UserID == nil || baseUserID != *req.UserID {
	// 	return nil, service_error.ErrTicketUnauthorizedAccess
	// }
	tempUserID := tokenKeys.GetUserID()
	tempAparID := tokenKeys.GetApartmentID()
	if tempAparID == constant.NilApartmentIDKeyToken && tokenKeys.GetRole() != entity.RoleAdmin {
		return nil, service_error.ErrTicketUnauthorizedAccess
	}
	ticket := &entity.Ticket{
		UserID:        &tempUserID,
		ApartmentID:   tempAparID,
		Title:         req.Title,
		Description:   req.Description,
		Body:          req.Body,
		Category:      req.Category,
		Accessibility: req.Accessibility,
		Status:        entity.TicketOpen,
	}
	err := s.repo.Create(ctx, ticket)
	if err != nil {
		return nil, err
	}

	var joinTags []entity.TicketAnnouncementTag
	var tags []entity.Tag
	if len(req.TagIDs) > 0 {
		tags, err = s.tagRepo.FindByIDs(ctx, req.TagIDs)
		if err != nil {
			return nil, err
		}

		joinTags = make([]entity.TicketAnnouncementTag, len(tags))
		for i, t := range tags {
			joinTags[i] = entity.TicketAnnouncementTag{
				TagID:    t.ID,
				TicketID: &ticket.ID,
				Tag:      t,
			}
		}
	}
	// fmt.Println(tags, joinTags, req.TagIDs)
	ticket.Tags = joinTags
	err = s.repo.CreateTags(ctx, ticket)
	if err != nil {
		return nil, err
	}

	return &dto.CreateTicketResponse{
		ID:          ticket.ID,
		UserID:      ticket.UserID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Body:        ticket.Body,
		Category:    string(ticket.Category),
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		Tags:        dto.MapTagsToSliceResponse(tags),
	}, nil
}

func (s *ticketService) Delete(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) error {

	ticket, err := s.getByIDSuperAccess(ctx, id)
	if err != nil {
		return err
	}

	isAdmin := tokenKeys.GetRole() == entity.RoleAdmin
	isManger := tokenKeys.GetRole() == entity.RoleManager
	isResident := tokenKeys.GetRole() == entity.RoleResident
	isNotForeignApartment := tokenKeys.GetApartmentID() == ticket.ApartmentID // TODO : WARNING!!!!!!
	isYourTicket := ticket.UserID != nil && tokenKeys.GetUserID() == *ticket.UserID

	cond_authz := isAdmin || (isManger && isNotForeignApartment) || (isResident && isNotForeignApartment)
	cond := cond_authz && isYourTicket
	if !cond {
		return service_error.ErrTicketUnauthorizedAccess
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

func (s *ticketService) GetByID(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) (*dto.TicketBaseResponse, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	isAdmin := tokenKeys.GetRole() == entity.RoleAdmin
	isManger := tokenKeys.GetRole() == entity.RoleManager
	isResident := tokenKeys.GetRole() == entity.RoleResident
	isNotForeignApartment := tokenKeys.GetApartmentID() == ticket.ApartmentID // TODO : WARNING!!!!!!
	isYourTicket := ticket.UserID != nil && tokenKeys.GetUserID() == *ticket.UserID
	isTicketPublic := ticket.Accessibility == entity.PublicTicket

	cond := isAdmin || (isManger && isNotForeignApartment) || (isResident && isNotForeignApartment && (isTicketPublic || isYourTicket))
	if !cond {
		return nil, service_error.ErrTicketUnauthorizedAccess
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

func (s *ticketService) GetByIDWithAllRelations(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) (*dto.TicketBaseResponseWithAllRelations, error) {
	ticket, err := s.repo.GetByIDWithAllRelations(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	isAdmin := tokenKeys.GetRole() == entity.RoleAdmin
	isManger := tokenKeys.GetRole() == entity.RoleManager
	isResident := tokenKeys.GetRole() == entity.RoleResident
	isNotForeignApartment := tokenKeys.GetApartmentID() == ticket.ApartmentID // TODO : WARNING!!!!!!
	isYourTicket := ticket.UserID != nil && tokenKeys.GetUserID() == *ticket.UserID
	isTicketPublic := ticket.Accessibility == entity.PublicTicket

	cond := isAdmin || (isManger && isNotForeignApartment) || (isResident && isNotForeignApartment && (isTicketPublic || isYourTicket))
	if !cond {
		return nil, service_error.ErrTicketUnauthorizedAccess
	}

	return &dto.TicketBaseResponseWithAllRelations{
		TicketBaseResponse: *dto.MapTicketToBaseResponse(ticket),
		Comments:           dto.MapCommentsToResponseWithUserSlice(ticket.Comments),
		User:               *dto.MapUserToUserResponse(&ticket.User),
		Apartment:          *dto.MapApartmentToResponse(&ticket.Apartment),
	}, nil
}

// TODO : This list it's not fully
func (s *ticketService) List(ctx context.Context, tokenKeys *dto.TokenKeys, filter dto.TicketFilterRequest) ([]dto.TicketBaseResponseWithCommentCount, error) {

	new_filter := domainRepo.TicketFilter{
		UserID:   filter.UserUUID,
		Status:   filter.Status,
		Category: filter.Category,
		Limit:    filter.Limit,
		Offset:   filter.Limit * (filter.Page - 1),
	}

	fl, err := s.repo.List(ctx, new_filter, tokenKeys.GetUserID(), tokenKeys.GetRole(), tokenKeys.GetApartmentID())
	if err != nil {
		return nil, err
	}

	return dto.MapTicketsToSliceResponseWithCount(fl), err
}

func (s *ticketService) Update(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID, req dto.UpdateTicketRequest) (*dto.TicketBaseResponse, error) {

	ticket, err := s.getByIDSuperAccess(ctx, id) // it must retrieve a record for checking after that
	if err != nil {
		return nil, err
	}

	isNotForeignApartment := tokenKeys.GetApartmentID() == ticket.ApartmentID // TODO : WARNING!!!!!!
	isYourTicket := ticket.UserID != nil && tokenKeys.GetUserID() == *ticket.UserID

	cond := isNotForeignApartment && isYourTicket
	if !cond {
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

func (s *ticketService) UpdateStatus(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID, status entity.TicketStatus) (*dto.TicketBaseResponse, error) {
	ticket, err := s.getByIDSuperAccess(ctx, id) // it must retrieve a record for checking after that
	if err != nil {
		return nil, err
	}

	isAdmin := tokenKeys.GetRole() == entity.RoleAdmin
	isManger := tokenKeys.GetRole() == entity.RoleManager
	isNotForeignApartment := tokenKeys.GetApartmentID() == ticket.ApartmentID // TODO : WARNING!!!!!!

	cond := isAdmin || (isManger && isNotForeignApartment)
	if !cond {
		return nil, service_error.ErrTicketUnauthorizedAccess
	}

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
