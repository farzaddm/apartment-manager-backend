package service

import (
	"apartment-manager-backend/internal/application/dto/user"
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

func (s *UserService) Update(ctx context.Context, user userdto.UpdateProfileRequest, id string) error {
	hashedPassword, err := s.hasher.Hash(user.Password)
	if err != nil {
		return errors.New("failed_to_hash_password")
	}
	user.Password = hashedPassword

	return s.UserRepo.Update(ctx, user, id)
}

func (s *UserService) Delete(ctx context.Context, id string) error {

	return s.UserRepo.Delete(ctx, id)
}

func (s *UserService) GetById(ctx context.Context, id string) (*userdto.UserProfileResponse, error) {
	rawUser, err := s.UserRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	profile := &userdto.UserProfileResponse{
		UserId:    rawUser.ID.String(),
		FirstName: rawUser.FirstName,
		LastName:  rawUser.LastName,
		Email:     rawUser.Email,
		Username:  rawUser.Username,
		Password:  rawUser.Password,
		Gender:    rawUser.Gender,
		Role:      rawUser.Role,
	}
	return profile, nil
}
