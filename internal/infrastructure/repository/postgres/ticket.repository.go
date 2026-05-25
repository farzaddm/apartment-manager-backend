package postgres

//TODO : Sort Res !!!
import (
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN tags ON tags.id = ticket_announcement_tags.tag_id"). //TODO : I'm not Sure Work Well or Not
				Order("tags.name ASC")
		}).
		Preload("Tags.Tag").
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
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("comments.committed_order ASC")
		}).
		Preload("Comments.User").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN tags ON tags.id = ticket_announcement_tags.tag_id"). //TODO : I'm not Sure Work Well or Not
				Order("tags.name ASC")
		}).
		Preload("Tags.Tag").
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
func (r *ticketRepository) List(ctx context.Context, filter domainRepo.TicketFilter, me uuid.UUID, role entity.UserRole) ([]domainRepo.TicketWithCommentCount, error) {
	var tickets []domainRepo.TicketWithCommentCount

	// TODO: Verify whether using Select here could cause issues.
	query := r.db.WithContext(ctx).
		Model(&entity.Ticket{}).
		Select("tickets.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.ticket_id = tickets.id").
		Order("tickets.created_at ASC").
		Group("tickets.id")

	if !(role == entity.RoleAdmin || role == entity.RoleManager) {
		query = query.Where(
			r.db.Where("tickets.accessibility = ?", "public").
				Or("tickets.user_id = ?", me),
		)
	}

	if filter.UserID != nil {
		nilUUID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

		if *filter.UserID == nilUUID {
			query = query.Where("tickets.user_id IS NULL")
		} else {
			query = query.Where("tickets.user_id = ?", *filter.UserID)
		}
	}

	if filter.Status != nil {
		query = query.Where("tickets.status = ?", *filter.Status)
	}
	if filter.Category != nil {
		query = query.Where("tickets.category = ?", *filter.Category)
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

func (r *ticketRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TicketStatus) (*entity.Ticket, error) {
	var ticket entity.Ticket

	result := r.db.WithContext(ctx).
		Model(&ticket).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Update("status", status)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &ticket, nil
}

func (r *ticketRepository) UpdateCategory(ctx context.Context, id uuid.UUID, category entity.TicketCategory) (*entity.Ticket, error) {
	var ticket entity.Ticket

	result := r.db.WithContext(ctx).
		Model(&ticket).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Update("category", category)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &ticket, nil
}

func (r *ticketRepository) UpdateContent(ctx context.Context, id uuid.UUID, title, description, body string) (*entity.Ticket, error) {
	var ticket entity.Ticket

	result := r.db.WithContext(ctx).
		Model(&ticket).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Updates(map[string]interface{}{
			"title":       title,
			"description": description,
			"body":        body,
		})

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &ticket, nil
}

func (r *ticketRepository) GetTicketsByUserId(ctx context.Context, userID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Tags").Find(&tickets).Error
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
