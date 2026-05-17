package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) domainRepo.TicketInterface {
	return &ticketRepository{db: db}
}

// /////////////////// Create / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////

func (r *ticketRepository) Create(ctx context.Context, ticket *entity.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

// /////////////////// Delete / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////

func (r *ticketRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Delete(&entity.Ticket{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// /////////////////// Get / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////

func (r *ticketRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
	var ticket entity.Ticket

	err := r.db.WithContext(ctx).
		Preload("Tags").
		First(&ticket, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &ticket, nil
}

func (r *ticketRepository) GetByIDWithAllRelations(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
	var ticket entity.Ticket

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Comments").
		Preload("Comments.User").
		Preload("Tags").
		First(&ticket, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &ticket, nil
}

// /////////////////// List / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////
// TODO: im not sure that is this code worked or not ? (maybe need some struct tags in domainRepo.TicketWithCommentCount)
func (r *ticketRepository) List(ctx context.Context, filter domainRepo.TicketFilter) ([]domainRepo.TicketWithCommentCount, error) {
	var tickets []domainRepo.TicketWithCommentCount

	// TODO: Verify whether using Select here could cause issues.
	query := r.db.WithContext(ctx).
		Model(&entity.Ticket{}).
		Select("tickets.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.ticket_id = tickets.id").
		Group("tickets.id")

	if filter.UserID != nil {
		query = query.Where("tickets.user_id = ?", *filter.UserID)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Scan(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

// /////////////////// Update / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////

func (r *ticketRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TicketStatus) error {
	result := r.db.WithContext(ctx).
		Model(&entity.Ticket{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *ticketRepository) UpdateCategory(ctx context.Context, id uuid.UUID, category entity.TicketCategory) error {
	result := r.db.WithContext(ctx).
		Model(&entity.Ticket{}).
		Where("id = ?", id).
		Update("category", category)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *ticketRepository) UpdateContent(ctx context.Context, id uuid.UUID, title, description, body string) error {
	result := r.db.WithContext(ctx).
		Model(&entity.Ticket{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"title":       title,
			"description": description,
			"body":        body,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *ticketRepository) GetTicketsByUserId(ctx context.Context, userID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Tags").Find(&tickets).Error
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
