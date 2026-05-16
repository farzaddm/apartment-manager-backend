package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/entity"
	domainJwt "apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"apartment-manager-backend/internal/domain/repository/redis"
	"apartment-manager-backend/pkg/hasher"
	"context"
	"errors"
	"time"
)

type RegisterService struct {
	userRepo     postgres.UserInterface
	otpRepo      redis.OTPInterFace
	refreshRepo  redis.RefreshTokenInterFace
	tokenService domainJwt.TokenServiceInterface
	hasher       *hasher.BcryptHasher
	refreshTTL   time.Duration
}

func NewRegisterService(
	userRepo postgres.UserInterface,
	otpRepo redis.OTPInterFace,
	refreshRepo redis.RefreshTokenInterFace,
	tokenService domainJwt.TokenServiceInterface,
	hasher *hasher.BcryptHasher,
) *RegisterService {
	return &RegisterService{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		refreshRepo:  refreshRepo,
		tokenService: tokenService,
		hasher:       hasher,
		refreshTTL:   7 * 24 * time.Hour,
	}
}

func (s *RegisterService) Execute(ctx context.Context, input dto.RegisterInput) (*dto.RegisterOutput, error) {
	isVerified, err := s.otpRepo.IsVerified(input.Phone)
	if err != nil || !isVerified {
		return nil, errors.New("phone_not_verified_or_expired")
	}

	emailExists, err := s.userRepo.ExistEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("email_already_exists")
	}

	existingUser, err := s.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username_already_exists")
	}

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, errors.New("failed_to_hash_password")
	}

	gender := entity.GenderFemale

	newUser := &entity.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  hashedPassword,
		Role:      "resident",
		Gender:    &gender,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	_ = s.otpRepo.DeleteVerified(input.Phone)

	accessToken, err := s.tokenService.GenerateAccessToken(newUser.ID, string(newUser.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(newUser.ID, string(newUser.Role))
	if err != nil {
		return nil, err
	}

	if s.refreshRepo != nil {
		if err := s.refreshRepo.Save(ctx, newUser.ID.String(), refreshToken, s.refreshTTL); err != nil {
			return nil, err
		}
	}

	return &dto.RegisterOutput{
		User:         newUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "registration_success",
	}, nil
}
