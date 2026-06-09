package service

import (
	"apartment-manager-backend/internal/application/dto"
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService interface {
	Create(ctx context.Context, tokenKeys *dto.TokenKeys, ticketID uuid.UUID, comment *dto.CreateCommentRequest) (*dto.CreateCommentResponse, error)

	Update(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID, comment *dto.UpdateCommentRequest) (*dto.CommentResponse, error)

	Delete(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) error

	GetByID(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) (*dto.CommentResponseWithAllRelations, error)

	GetLastOrderByTicketID(ctx context.Context, tokenKeys *dto.TokenKeys, ticketID uuid.UUID) (*int, error)
}

type commentService struct {
	commentRepo domainRepo.CommentInterface
	ticketRepo  domainRepo.TicketInterface
}

func NewCommentService(commentRepo domainRepo.CommentInterface, ticketRepo domainRepo.TicketInterface) CommentService {
	return &commentService{commentRepo: commentRepo, ticketRepo: ticketRepo}
}

// TODO : for its handler
func (s *commentService) Create(ctx context.Context, tokenKeys *dto.TokenKeys, ticketID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CreateCommentResponse, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
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
	isTicketPublic := ticket.Accessibility == entity.PublicTicket
	isYourTicket := ticket.UserID != nil && tokenKeys.GetUserID() == *ticket.UserID
	isNotClosedTicket := ticket.Status != entity.TicketClosed

	cond_authz := isAdmin || (isManger && isNotForeignApartment) || (isResident && isNotForeignApartment && (isTicketPublic || isYourTicket))
	cond := cond_authz && isNotClosedTicket

	if !cond {
		return nil, service_error.ErrCommentUnauthorizedAccess
	}

	lastOrder, err := s.commentRepo.GetLastOrderByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if lastOrder == nil {
		return nil, service_error.ErrCommentNotFound
	}

	// order := 1
	// if lastOrder != nil {
	// 	order = *lastOrder + 1
	// }
	temp := tokenKeys.GetUserID()
	comment := &entity.Comment{
		TicketID:       ticketID,
		Body:           req.Body,
		CommittedOrder: *lastOrder + 1,
		UserID:         &temp,
	}
	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}
	return &dto.CreateCommentResponse{
		ID:             comment.ID,
		UserID:         comment.UserID,
		TicketID:       comment.TicketID,
		Body:           comment.Body,
		CommittedOrder: comment.CommittedOrder,
		CreatedAt:      comment.CreatedAt,
	}, nil
}

func (s *commentService) Update(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error) {
	comm, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if comm == nil {
		return nil, service_error.ErrTicketOfCommentOrCommentNotFound
	}

	ticket, err := s.ticketRepo.GetByID(ctx, comm.TicketID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	isNotForeignApartment := tokenKeys.GetApartmentID() == ticket.ApartmentID // TODO : WARNING!!!!!!
	isYourComment := comm.UserID != nil && tokenKeys.GetUserID() == *comm.UserID
	isNotClosedTicket := ticket.Status != entity.TicketClosed

	cond_authz := isNotForeignApartment && isYourComment
	cond := cond_authz && isNotClosedTicket
	if !cond {
		return nil, service_error.ErrCommentUnauthorizedAccess
	}

	comment := &entity.Comment{
		Body: req.Body,
	}
	comm, err = s.commentRepo.Update(ctx, id, comment)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrCommentNotFound
		}
		return nil, err
	}
	return dto.MapCommentToResponse(comm), nil
}

func (s *commentService) Delete(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) error {
	comm, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if comm == nil {
		return service_error.ErrCommentNotFound
	}

	ticket, err := s.ticketRepo.GetByID(ctx, comm.TicketID)
	if err != nil {
		return err
	}
	if ticket == nil {
		return service_error.ErrTicketNotFound
	}

	isAdmin := tokenKeys.GetRole() == entity.RoleAdmin
	isResident := tokenKeys.GetRole() == entity.RoleResident
	isManger := tokenKeys.GetRole() == entity.RoleManager
	isNotForeignApartment := tokenKeys.GetApartmentID() == ticket.ApartmentID // TODO : WARNING!!!!!!
	isYourComment := comm.UserID != nil && tokenKeys.GetUserID() == *comm.UserID
	isNotClosedTicket := ticket.Status != entity.TicketClosed

	cond_authz := isAdmin || (isManger && isNotForeignApartment) || (isResident && isNotForeignApartment && isYourComment)
	cond := cond_authz && isNotClosedTicket
	if !cond {
		return service_error.ErrCommentUnauthorizedAccess
	}

	err = s.commentRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrCommentNotFound
		}
		return err
	}
	return nil
}

func (s *commentService) GetByID(ctx context.Context, tokenKeys *dto.TokenKeys, id uuid.UUID) (*dto.CommentResponseWithAllRelations, error) {
	comm, err := s.commentRepo.GetWithTicket(ctx, id)
	if err != nil {
		return nil, err
	}
	if comm == nil {
		return nil, service_error.ErrTicketOfCommentOrCommentNotFound
	}

	ticket, err := s.ticketRepo.GetByID(ctx, comm.TicketID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	isAdmin := tokenKeys.GetRole() == entity.RoleAdmin
	isResident := tokenKeys.GetRole() == entity.RoleResident
	isManger := tokenKeys.GetRole() == entity.RoleManager
	isNotForeignApartment := tokenKeys.GetApartmentID() == ticket.ApartmentID // TODO : WARNING!!!!!!
	isTicketPublic := ticket.Accessibility == entity.PublicTicket
	isYourTicket := ticket.UserID != nil && tokenKeys.GetUserID() == *ticket.UserID

	cond := isAdmin || (isManger && isNotForeignApartment) || (isResident && isNotForeignApartment && (isTicketPublic || isYourTicket))

	if !cond {
		return nil, service_error.ErrCommentUnauthorizedAccess
	}

	c, err := s.commentRepo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, service_error.ErrCommentNotFound
	}
	return &dto.CommentResponseWithAllRelations{
		CommentResponse: *dto.MapCommentToResponse(c),
		User:            *dto.MapUserToUserResponse(&c.User),
		Ticket:          *dto.MapTicketToBaseResponse(&c.Ticket),
	}, nil
}

func (s *commentService) GetLastOrderByTicketID(ctx context.Context, tokenKeys *dto.TokenKeys, ticketID uuid.UUID) (*int, error) {
	comm, err := s.commentRepo.GetWithTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if comm == nil {
		return nil, service_error.ErrTicketOfCommentOrCommentNotFound
	}

	ticket, err := s.ticketRepo.GetByID(ctx, comm.TicketID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, service_error.ErrTicketNotFound
	}

	isAdmin := tokenKeys.GetRole() == entity.RoleAdmin
	isResident := tokenKeys.GetRole() == entity.RoleResident
	isManger := tokenKeys.GetRole() == entity.RoleManager
	isNotForeignApartment := tokenKeys.GetApartmentID() == ticket.ApartmentID // TODO : WARNING!!!!!!
	isTicketPublic := ticket.Accessibility == entity.PublicTicket
	isYourTicket := ticket.UserID != nil && tokenKeys.GetUserID() == *ticket.UserID

	cond := isAdmin || (isManger && isNotForeignApartment) || (isResident && isNotForeignApartment && (isTicketPublic || isYourTicket))
	if !cond {
		return nil, service_error.ErrCommentUnauthorizedAccess
	}

	count, err := s.commentRepo.GetLastOrderByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if count == nil {
		return nil, service_error.ErrTicketOfCommentNotFound
	}
	return count, nil
}
