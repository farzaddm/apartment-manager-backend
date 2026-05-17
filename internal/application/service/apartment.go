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

type ApartmentService interface {
	Create(ctx context.Context, apartment *entity.Apartment) error
	Update(ctx context.Context, apartment *entity.Apartment) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)
	GetByIDWithRelations(ctx context.Context, id uuid.UUID, relations ...string) (*entity.Apartment, error)

	GetWithUsers(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)
	GetWithAnnouncements(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)
	GetWithRules(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)
	GetWithInviteCodes(ctx context.Context, id uuid.UUID) (*entity.Apartment, error)

	// List(ctx context.Context) ([]entity.Apartment, error)

	// ListWithRelations(ctx context.Context, relations ...string) ([]entity.Apartment, error)

	// ListWithUsers(ctx context.Context) ([]entity.Apartment, error)
	// ListWithAnnouncements(ctx context.Context) ([]entity.Apartment, error)
	// ListWithRules(ctx context.Context) ([]entity.Apartment, error)
	// ListWithInviteCodes(ctx context.Context) ([]entity.Apartment, error)
}

type apartmentService struct {
	apartmentRepo domainRepo.ApartmentInterface
}

func NewApartmentService(apartmentRepo domainRepo.ApartmentInterface) ApartmentService {
	return &apartmentService{
		apartmentRepo: apartmentRepo,
	}
}

func (s *apartmentService) Create(ctx context.Context, apartment *entity.Apartment) error {
	return s.apartmentRepo.Create(ctx, apartment)
}

func (s *apartmentService) Update(ctx context.Context, apartment *entity.Apartment) error {
	err := s.apartmentRepo.Update(ctx, apartment)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrApartmentNotFound
		}
		return err
	}

	return nil
}

func (s *apartmentService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.apartmentRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrApartmentNotFound
		}
		return err
	}

	return nil
}

func (s *apartmentService) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := s.apartmentRepo.Exists(ctx, id)
	if err != nil {
		return false, err
	}
	return *exists, nil
}

func (s *apartmentService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	apartment, err := s.apartmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if apartment == nil {
		return nil, service_error.ErrApartmentNotFound
	}

	return apartment, nil
}

func (s *apartmentService) GetByIDWithRelations(ctx context.Context, id uuid.UUID, relations ...string) (*entity.Apartment, error) {
	apartment, err := s.apartmentRepo.GetByIDWithRelations(ctx, id, relations...)
	if err != nil {
		return nil, err
	}

	if apartment == nil {
		return nil, service_error.ErrApartmentNotFound
	}

	return apartment, nil
}

func (s *apartmentService) GetWithUsers(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	apartment, err := s.apartmentRepo.GetWithUsers(ctx, id)
	if err != nil {
		return nil, err
	}

	if apartment == nil {
		return nil, service_error.ErrApartmentNotFound
	}

	return apartment, nil
}

func (s *apartmentService) GetWithAnnouncements(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	apartment, err := s.apartmentRepo.GetWithAnnouncements(ctx, id)
	if err != nil {
		return nil, err
	}

	if apartment == nil {
		return nil, service_error.ErrApartmentNotFound
	}

	return apartment, nil
}

func (s *apartmentService) GetWithRules(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	apartment, err := s.apartmentRepo.GetWithRules(ctx, id)
	if err != nil {
		return nil, err
	}

	if apartment == nil {
		return nil, service_error.ErrApartmentNotFound
	}

	return apartment, nil
}

func (s *apartmentService) GetWithInviteCodes(ctx context.Context, id uuid.UUID) (*entity.Apartment, error) {
	apartment, err := s.apartmentRepo.GetWithInviteCodes(ctx, id)
	if err != nil {
		return nil, err
	}

	if apartment == nil {
		return nil, service_error.ErrApartmentNotFound
	}

	return apartment, nil
}

// func (s *apartmentService) List(ctx context.Context) ([]entity.Apartment, error) {
// 	return s.apartmentRepo.List(ctx)
// }

// func (s *apartmentService) ListWithRelations(ctx context.Context, relations ...string) ([]entity.Apartment, error) {
// 	return s.apartmentRepo.ListWithRelations(ctx, relations...)
// }

// func (s *apartmentService) ListWithUsers(ctx context.Context) ([]entity.Apartment, error) {
// 	return s.apartmentRepo.ListWithUsers(ctx)
// }

// func (s *apartmentService) ListWithAnnouncements(ctx context.Context) ([]entity.Apartment, error) {
// 	return s.apartmentRepo.ListWithAnnouncements(ctx)
// }

// func (s *apartmentService) ListWithRules(ctx context.Context) ([]entity.Apartment, error) {
// 	return s.apartmentRepo.ListWithRules(ctx)
// }

// func (s *apartmentService) ListWithInviteCodes(ctx context.Context) ([]entity.Apartment, error) {
// 	return s.apartmentRepo.ListWithInviteCodes(ctx)
// }
