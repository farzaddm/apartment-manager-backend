package service

import (
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService interface {
	Create(ctx context.Context, comment *entity.Comment) error

	Update(ctx context.Context, comment *entity.Comment) error

	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error)

	ListByTicketID(ctx context.Context, ticketID uuid.UUID) ([]entity.Comment, error)

	GetLastOrderByTicketID(ctx context.Context, ticketID uuid.UUID) (*int, error)
}

type commentService struct {
	commentRepo domainRepo.CommentInterface
}

func NewCommentService(commentRepo domainRepo.CommentInterface) CommentService {
	return &commentService{commentRepo: commentRepo}
}

func (s *commentService) Create(ctx context.Context, comment *entity.Comment) error {
	lastOrder, err := s.commentRepo.GetLastOrderByTicketID(comment.TicketID)
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

	comment.CommittedOrder = *lastOrder + 1

	return s.commentRepo.Create(comment)
}

// func (s *commentService) Update(ctx context.Context, comment *dto.UpdateCommentRequest) error {
// 	err := s.commentRepo.Update(comment)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return service_error.ErrCommentNotFound
// 		}
// 		return err
// 	}
// 	return nil
// }

func (s *commentService) Delete(ctx context.Context, id uuid.UUID, baseUserID uuid.UUID, role entity.UserRole) error {
	c, err := s.GetByID(id, baseUserID, role)
	if err != nil {
		return err
	}
	err = s.commentRepo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrCommentNotFound
		}
		return err
	}
	return nil
}

func (s *commentService) GetByID(ctx context.Context, id uuid.UUID, baseUserID uuid.UUID, role entity.UserRole) (*entity.Comment, error) {
	c, err := s.commentRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrCommentNotFound
		}
		return nil, err
	}
	if c.UserID == nil || *c.UserID != baseUserID {
		if !(role == entity.RoleAdmin || role == entity.RoleManager) {
			return nil, service_error.ErrCommentUnauthorizedAccess
		}
	}
	return c, nil
}

func (s *commentService) ListByTicketID(ctx context.Context, ticketID uuid.UUID, baseUserID uuid.UUID, role entity.UserRole) ([]entity.Comment, error) {
	l, err := s.commentRepo.ListByTicketID(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrCommentNotFound
		}
		return nil, err
	}
	new_l := make([]entity.Comment, 0)
	for i := range l {
		if l[i].UserID == nil || *l[i].UserID != baseUserID {
			if !(role == entity.RoleAdmin || role == entity.RoleManager) {
				continue
			}
		}
		new_l = append(new_l, l[i])
	}
	return l, nil
}

func (s *commentService) GetLastOrderByTicketID(ctx context.Context, ticketID uuid.UUID) (*int, error) {
	count, err := s.commentRepo.GetLastOrderByTicketID(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrCommentNotFound
		}
		return nil, err
	}
	return count, nil
}
