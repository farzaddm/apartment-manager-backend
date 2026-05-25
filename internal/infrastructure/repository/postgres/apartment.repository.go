package postgres

//TODO : Sort Lists!
import (
	"apartment-manager-backend/internal/domain/entity"
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type apartmentRepository struct {
	db *gorm.DB
}

func NewApartmentRepository(db *gorm.DB) domainRepo.ApartmentInterface {
	return &apartmentRepository{db: db}
}

// /////////////////// Create / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////

func (r *apartmentRepository) Create(ctx context.Context, apartment *entity.Apartment) error {
	return r.db.WithContext(ctx).Create(apartment).Error
}

// /////////////////// Update / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////

func (r *apartmentRepository) Update(ctx context.Context, id uuid.UUID, apartment *entity.Apartment) (*entity.Apartment, error) {
	result := r.db.WithContext(ctx).
		Model(apartment).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(apartment)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return apartment, nil
}

// /////////////////// Delete / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////

func (r *apartmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Delete(&entity.Apartment{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// /////////////////// EXIST / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////

func (r *apartmentRepository) Exists(ctx context.Context, id uuid.UUID) (*bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&entity.Apartment{}).
		Where("id = ?", id).
		Count(&count).Error

	if err != nil {
		return nil, err
	}
	b := count > 0
	return &b, nil
}

// /////////////////// Get / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////

func (r *apartmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	var apartment entity.Apartment

	err := r.db.WithContext(ctx).
		First(&apartment, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apartment, nil
}

// func (r *apartmentRepository) GetByIDWithRelations(ctx context.Context, id uuid.UUID, relations ...string) (*entity.Apartment, error) {
// 	var apartment entity.Apartment

// 	query := r.db.WithContext(ctx)

// 	for _, relation := range relations {
// 		query = query.Preload(relation)
// 	}

// 	err := query.
// 		First(&apartment, "id = ?", id).
// 		Error

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &apartment, nil
// }

func (r *apartmentRepository) GetWithUsers(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	var apartment entity.Apartment

	err := r.db.WithContext(ctx).
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Order("users.created_at ASC")
		}).
		First(&apartment, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apartment, nil
}

func (r *apartmentRepository) GetWithAnnouncements(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	var apartment entity.Apartment

	err := r.db.WithContext(ctx).
		Preload("Announcements", func(db *gorm.DB) *gorm.DB {
			return db.Order("is_pinned DESC").
				Order(`CASE 
                WHEN announcements.order = 'warning' THEN 1 
                WHEN announcements.order = 'very_important' THEN 2 
                WHEN announcements.order = 'important' THEN 3 
                WHEN announcements.order = 'other' THEN 4 
                ELSE 5 
            END ASC`)
		}).
		Preload("Announcements.Tags", func(db *gorm.DB) *gorm.DB {
			return db.Joins("INNER JOIN tags ON tags.id = ticket_announcement_tags.tag_id"). //TODO : I'm not Sure Work Well or Not
														Order("tags.name ASC")
		}).
		Preload("Announcements.Tags.Tag").
		First(&apartment, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apartment, nil
}

func (r *apartmentRepository) GetWithRules(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	var apartment entity.Apartment

	err := r.db.WithContext(ctx).
		Preload("Rules", func(db *gorm.DB) *gorm.DB {
			return db.Order("rules.created_at ASC")
		}).
		First(&apartment, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apartment, nil
}

func (r *apartmentRepository) GetWithInviteCodes(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	var apartment entity.Apartment

	err := r.db.WithContext(ctx).
		Preload("InviteCodes", func(db *gorm.DB) *gorm.DB {
			return db.Order("invite_codes.expire_at ASC")
		}).
		First(&apartment, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apartment, nil
}

// /////////////////// LIST / //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// //////////////////// ///////////////////
func (r *apartmentRepository) List(ctx context.Context) ([]entity.Apartment, error) {
	var apartments []entity.Apartment

	err := r.db.WithContext(ctx).
		Find(&apartments).
		Order("create_at DESC").
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return apartments, nil
}

// func (r *apartmentRepository) ListWithRelations(ctx context.Context, relations ...string) ([]entity.Apartment, error) {
// 	var apartments []entity.Apartment

// 	query := r.db.WithContext(ctx)

// 	for _, relation := range relations {
// 		query = query.Preload(relation)
// 	}

// 	err := query.
// 		Find(&apartments).
// 		Error

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return apartments, nil
// }

func (r *apartmentRepository) ListWithUsers(ctx context.Context) ([]entity.Apartment, error) {
	var apartments []entity.Apartment

	err := r.db.WithContext(ctx).
		Preload("Users").
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Order("users.created_at DESC")
		}).
		Find(&apartments).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return apartments, nil
}

func (r *apartmentRepository) ListWithAnnouncements(ctx context.Context) ([]entity.Apartment, error) {
	var apartments []entity.Apartment

	err := r.db.WithContext(ctx).
		Preload("Announcements", func(db *gorm.DB) *gorm.DB {
			return db.Order("is_pinned DESC").
				Order(`CASE 
                WHEN announcements.order = 'warning' THEN 1 
                WHEN announcements.order = 'very_important' THEN 2 
                WHEN announcements.order = 'important' THEN 3 
                WHEN announcements.order = 'other' THEN 4 
                ELSE 5 
            END ASC`)
		}).
		Preload("Announcements.Tags", func(db *gorm.DB) *gorm.DB {
			return db.Joins("INNER JOIN tags ON tags.id = ticket_announcement_tags.tag_id"). //TODO : I'm not Sure Work Well or Not
														Order("tags.name ASC")
		}).
		Preload("Announcements.Tags.Tag").
		Find(&apartments).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return apartments, nil
}

func (r *apartmentRepository) ListWithRules(ctx context.Context) ([]entity.Apartment, error) {
	var apartments []entity.Apartment

	err := r.db.WithContext(ctx).
		Preload("Rules", func(db *gorm.DB) *gorm.DB {
			return db.Order("rules.created_at ASC")
		}).
		Find(&apartments).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return apartments, nil
}

func (r *apartmentRepository) ListWithInviteCodes(ctx context.Context) ([]entity.Apartment, error) {
	var apartments []entity.Apartment

	err := r.db.WithContext(ctx).
		Preload("InviteCodes", func(db *gorm.DB) *gorm.DB {
			return db.Order("invite_codes.expire_at ASC")
		}).
		Find(&apartments).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return apartments, nil
}

func (r *apartmentRepository) GetApartmentManagerID(ctx context.Context, apartmentID string) (string, error) {
	var manager entity.User

	err := r.db.WithContext(ctx).
		Where("apartment_id = ? AND role = ?", apartmentID, "manager").
		First(&manager).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("no manager found for this building")
		}
		return "", err
	}

	return manager.ID.String(), nil
}
