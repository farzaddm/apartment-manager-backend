package postgres

//TODO : Sort Res List!
import (
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) domainRepo.CommentInterface {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) Update(ctx context.Context, id uuid.UUID, comment *entity.Comment) (*entity.Comment, error) {
	result := r.db.WithContext(ctx).
		Model(comment).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Updates(comment)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// حالا اینجا comment کاملاً پر شده (فیلدهای قدیمی + فیلدهای جدید)
	return comment, nil
}
func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.Comment{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func (r *commentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.WithContext(ctx).Preload("User").Preload("Ticket").First(&comment, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) ListByTicketID(ctx context.Context, ticketID uuid.UUID) ([]entity.Comment, error) {
	var comments []entity.Comment
	err := r.db.WithContext(ctx).
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

func (r *commentRepository) GetLastOrderByTicketID(ctx context.Context, ticketID uuid.UUID) (*int, error) {
	var lastOrder int
	err := r.db.WithContext(ctx).
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

func (r *commentRepository) GetWithUser(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.WithContext(ctx).Preload("User").First(&comment, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, err
}

func (r *commentRepository) GetWithTicket(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.WithContext(ctx).Preload("Ticket").First(&comment, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, err
}

func (r *commentRepository) GetWithAllRelations(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.WithContext(ctx).Preload("User").Preload("Ticket").First(&comment, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, err
}
