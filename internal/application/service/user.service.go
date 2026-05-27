package service

import (
	domainRepo "apartment-manager-backend/internal/domain/repository/postgres"

	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"apartment-manager-backend/pkg/hasher"
	"context"
	"errors"
)

type UserService struct {
	UserRepo postgres.UserInterface
	hasher   *hasher.BcryptHasher
}

func NewUserService(userRepo postgres.UserInterface, hasher *hasher.BcryptHasher) *UserService {
	return &UserService{UserRepo: userRepo, hasher: hasher}
}

func (s *UserService) Update(ctx context.Context, user dto.UpdateProfileRequest, id string) error {
	d_user := domainRepo.UpdateProfileRequest{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Gender:    user.Gender,
	}
	return s.UserRepo.Update(ctx, d_user, id)
}

func (s *UserService) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest, id string) error {
	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		return errors.New("failed_to_hash_password")
	}
	return s.UserRepo.ChangePassword(ctx, hashedPassword, id)
}

func (s *UserService) Delete(ctx context.Context, id string) error {

	return s.UserRepo.Delete(ctx, id)
}

func (s *UserService) GetById(ctx context.Context, id string) (*dto.UserProfileResponse, error) {
	rawUser, unitNumber, err := s.UserRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if rawUser == nil {
		return nil, nil
	}

	profile := &dto.UserProfileResponse{
		ID:              rawUser.ID,
		CreatedAt:       rawUser.CreatedAt,
		ApartmentID:     rawUser.ApartmentID,
		UnitNumber:      unitNumber,
		FirstName:       rawUser.FirstName,
		LastName:        rawUser.LastName,
		Username:        rawUser.Username,
		Email:           rawUser.Email,
		Phone:           rawUser.Phone,
		Role:            string(rawUser.Role),
		Gender:          string(*rawUser.Gender),
		ProfileImageURL: rawUser.ProfileImageURL,
	}

	return profile, nil
}

func (s *UserService) SetProfileImage(ctx context.Context, userID string, filePath string) error {
	return s.UserRepo.UpdateProfileImage(ctx, userID, filePath)
}
