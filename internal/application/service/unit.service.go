package service

import (
	"apartment-manager-backend/internal/application/dto"
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UnitService interface {
	Create(ctx context.Context, apartment_id uuid.UUID, req *dto.CreateUnitRequest) (*dto.CreateUnitResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateUnitRequest) (*dto.UnitResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	PopUser(ctx context.Context, id uuid.UUID) (*dto.UnitResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.UnitResponse, error)
	PushUser(ctx context.Context, id uuid.UUID, req *dto.PushUserToUnitRequest) (*dto.UnitResponse, error)
}

type unitService struct {
	repo     postgres.UnitInterface
	userRepo postgres.UserInterface
}

func NewUnitService(repo postgres.UnitInterface, userRepo postgres.UserInterface) UnitService {
	return &unitService{repo: repo, userRepo: userRepo}
}

func (s *unitService) Create(ctx context.Context, apartment_id uuid.UUID, req *dto.CreateUnitRequest) (*dto.CreateUnitResponse, error) {
	unit := &entity.Unit{
		ApartmentID: apartment_id,
		UnitNumber:  req.UnitNumber,
		Floor:       req.Floor,
	}

	err := s.repo.Create(ctx, unit)
	if err != nil {
		return nil, err
	}
	return &dto.CreateUnitResponse{
		ID:          unit.ID,
		ApartmentID: unit.ApartmentID,
		UserID:      unit.UserID,
		UnitNumber:  unit.UnitNumber,
		Floor:       unit.Floor,
		CreatedAt:   unit.CreatedAt,
	}, nil
}
func (s *unitService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateUnitRequest) (*dto.UnitResponse, error) {
	existingUnit, err := s.repo.GetByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrUnitNotFound
		}
		return nil, err
	}

	existingUnit.UnitNumber = req.UnitNumber
	existingUnit.Floor = req.Floor

	updatedUnit, err := s.repo.Update(ctx, existingUnit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrUnitNotFound
		}
		return nil, service_error.ErrUnitUpdateFailed
	}

	return dto.MapUnitToResponse(updatedUnit), nil
}

func (s *unitService) PopUser(ctx context.Context, id uuid.UUID) (*dto.UnitResponse, error) {
	updatedUnit, err := s.repo.PopUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrUnitNotFound
		}
		return nil, err
	}

	return dto.MapUnitToResponse(updatedUnit), nil
}

func (s *unitService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service_error.ErrUnitNotFound
		}
		return err
	}

	return nil
}

func (s *unitService) GetByID(ctx context.Context, id uuid.UUID) (*dto.UnitResponse, error) {
	unit, err := s.repo.GetByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrUnitNotFound
		}
		return nil, err
	}
	return dto.MapUnitToResponse(unit), nil
}

func (s *unitService) PushUser(ctx context.Context, id uuid.UUID, req *dto.PushUserToUnitRequest) (*dto.UnitResponse, error) {
	updatedUnit, err := s.repo.PushUser(ctx, id, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrUnitNotFound
		}
		return nil, err
	}
	data := dto.MapUnitToResponse(updatedUnit)
	user, _, err := s.userRepo.GetById(ctx, req.UserID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service_error.ErrUserNotFound
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	data.User = dto.MapUserToUserResponse(user)
	return data, nil
}
