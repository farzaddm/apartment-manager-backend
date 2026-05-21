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
	Create(ctx context.Context, ticketID uuid.UUID, comment *dto.CreateCommentRequest) (*entity.Comment, error)

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

func (s *commentService) Create(ctx context.Context, ticketID uuid.UUID, req *dto.CreateCommentRequest) (*entity.Comment, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
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

	if ticket.Accessability == entity.PrivateTicket && role == entity.RoleResident && (ticket.UserID == nil || baseUserID != *ticket.UserID) {
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

	comment := &entity.Comment{
		TicketID:       ticketID,
		Body:           req.Body,
		CommittedOrder: *lastOrder + 1,
		UserID:         &baseUserID,
	}
	err = s.commentRepo.Create(ctx, comment)
	return comment, nil
}

func (s *commentService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateCommentRequest) error {
	comm, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if comm == nil {
		return service_error.ErrTicketOfCommentOrCommentNotFound
	}

	rawBaseUserID := ctx.Value("user_id") // IT MUST BE EXIST!
	if rawBaseUserID == nil {
		return service_error.ErrUserIDNotFoundInContext
	}
	str_baseUserID := rawBaseUserID.(string)
	baseUserID, err := uuid.Parse(str_baseUserID)

	if err != nil {
		return service_error.ErrCommonParseStrToUUID
	}

	if comm.UserID == nil || *comm.UserID != baseUserID {
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

	if (comm.UserID == nil || *comm.UserID != baseUserID) && role == entity.RoleResident {
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
	if comm.Ticket.Accessability == entity.PrivateTicket && (comm.Ticket.UserID == nil || *comm.Ticket.UserID != baseUserID) && role == entity.RoleResident {
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
	if comm.Ticket.Accessability == entity.PrivateTicket && (comm.Ticket.UserID == nil || *comm.Ticket.UserID != baseUserID) && role == entity.RoleResident {
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
