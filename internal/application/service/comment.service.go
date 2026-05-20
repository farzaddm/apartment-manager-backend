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
	Create(ctx context.Context, comment *dto.CreateCommentRequest) error

	Update(ctx context.Context, id uuid.UUID, comment *dto.UpdateCommentRequest) error

	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error)

	GetLastOrderByTicketID(ctx context.Context, ticketID uuid.UUID) (*int, error)
}

type commentService struct {
	commentRepo domainRepo.CommentInterface
	ticketRepo  domainRepo.TicketInterface
}

func NewCommentService(commentRepo domainRepo.CommentInterface, ticketRepo domainRepo.TicketInterface) CommentService {
	return &commentService{commentRepo: commentRepo, ticketRepo: ticketRepo}
}

func (s *commentService) Create(ctx context.Context, req *dto.CreateCommentRequest) error {
	ticket, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrTicketNotFound
		}
		return err
	}

	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return service_error.ErrUserIDNotFoundInContext
	}
	role := ctx.Value("role") // IT MUST BE EXIST!
	if role == nil {
		return service_error.ErrUserRoleNotFoundInContext
	}
	if ticket.Accessability == entity.PrivateTicket && role == entity.RoleResident {
		return service_error.ErrCommentUnauthorizedAccess
	}

	lastOrder, err := s.commentRepo.GetLastOrderByTicketID(ctx, req.TicketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrCommentNotFound
		}
		return err
	}

	// order := 1
	// if lastOrder != nil {
	// 	order = *lastOrder + 1
	// }

	cr_comment := &entity.Comment{
		TicketID:       req.TicketID,
		Body:           req.Body,
		CommittedOrder: *lastOrder + 1,
	}
	return s.commentRepo.Create(ctx, cr_comment)
}

func (s *commentService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateCommentRequest) error {
	comm, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if comm == nil {
		return service_error.ErrTicketOfCommentOrCommentNotFound
	}

	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return service_error.ErrUserIDNotFoundInContext
	}
	role := ctx.Value("role") // IT MUST BE EXIST!
	if role == nil {
		return service_error.ErrUserRoleNotFoundInContext
	}
	if comm.UserID == nil || *comm.UserID == baseUserID {
		return service_error.ErrCommentUnauthorizedAccess
	}
	comment := &entity.Comment{
		Body: req.Body,
	}
	err = s.commentRepo.Update(ctx, id, comment)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrCommentNotFound
		}
		return err
	}
	return nil
}

func (s *commentService) Delete(ctx context.Context, id uuid.UUID) error {
	comm, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if comm == nil {
		return service_error.ErrCommentNotFound
	}

	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return service_error.ErrUserIDNotFoundInContext
	}
	role := ctx.Value("role") // IT MUST BE EXIST!
	if role == nil {
		return service_error.ErrUserRoleNotFoundInContext
	}
	if (comm.UserID == nil || *comm.UserID == baseUserID) && role == entity.RoleResident {
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

func (s *commentService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	comm, err := s.commentRepo.GetWithTicket(ctx, id)
	if err != nil {
		return nil, err
	}
	if comm == nil {
		return nil, service_error.ErrTicketOfCommentOrCommentNotFound
	}
	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	role := ctx.Value("role") // IT MUST BE EXIST!
	if role == nil {
		return nil, service_error.ErrUserRoleNotFoundInContext
	}
	if comm.Ticket.Accessability == entity.PrivateTicket && (comm.Ticket.UserID == nil || *comm.Ticket.UserID == baseUserID) && role == entity.RoleResident {
		return nil, service_error.ErrCommentUnauthorizedAccess
	}

	c, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, service_error.ErrCommentNotFound
	}
	return c, nil
}

func (s *commentService) GetLastOrderByTicketID(ctx context.Context, ticketID uuid.UUID) (*int, error) {
	comm, err := s.commentRepo.GetWithTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if comm == nil {
		return nil, service_error.ErrTicketOfCommentOrCommentNotFound
	}
	baseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if baseUserID == nil {
		return nil, service_error.ErrUserIDNotFoundInContext
	}
	role := ctx.Value("role") // IT MUST BE EXIST!
	if role == nil {
		return nil, service_error.ErrUserRoleNotFoundInContext
	}
	if comm.Ticket.Accessability == entity.PrivateTicket && (comm.Ticket.UserID == nil || *comm.Ticket.UserID == baseUserID) && role == entity.RoleResident {
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
