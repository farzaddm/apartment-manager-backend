package service

import (
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService interface {
}

type commentService struct {
	commentRepo domainRepo.CommentInterface
}

func NewCommentService(commentRepo domainRepo.CommentInterface) CommentService {
	return &commentService{commentRepo: commentRepo}
}

func (s *commentService) Create(comment *entity.Comment) error {
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

func (s *commentService) Update(comment *entity.Comment) error {
	err := s.commentRepo.Update(comment)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrCommentNotFound
		}
		return err
	}
	return nil
}

func (s *commentService) Delete(id uuid.UUID) error {
	err := s.commentRepo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrCommentNotFound
		}
		return err
	}
	return nil
}

func (s *commentService) GetByID(id uuid.UUID) (*entity.Comment, error) {
	c, err := s.commentRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrCommentNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *commentService) ListByTicketID(ticketID uuid.UUID) ([]entity.Comment, error) {
	l, err := s.commentRepo.ListByTicketID(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrCommentNotFound
		}
		return nil, err
	}
	return l, nil
}

func (s *commentService) GetLastOrderByTicketID(ticketID uuid.UUID) (*int, error) {
	count, err := s.commentRepo.GetLastOrderByTicketID(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrCommentNotFound
		}
		return nil, err
	}
	return count, nil
}
