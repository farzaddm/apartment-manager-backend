package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) domainRepo.CommentInterface {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *entity.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) Update(comment *entity.Comment) error {
	result := r.db.Save(comment)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *commentRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&entity.Comment{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func (r *commentRepository) GetByID(id uuid.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.Preload("User").Preload("Ticket").First(&comment, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) ListByTicketID(ticketID uuid.UUID) ([]entity.Comment, error) {
	var comments []entity.Comment
	err := r.db.
		Where("ticket_id = ?", ticketID).
		Order("committed_order ASC").
		Preload("User").
		Find(&comments).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return comments, err
}

func (r *commentRepository) GetLastOrderByTicketID(ticketID uuid.UUID) (*int, error) {
	var lastOrder int
	err := r.db.
		Model(&entity.Comment{}).
		Where("ticket_id = ?", ticketID).
		Select("COALESCE(MAX(committed_order), 0)").
		Scan(&lastOrder).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &lastOrder, err
}
