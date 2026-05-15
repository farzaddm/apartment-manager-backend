package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"apartment-manager-backend/internal/domain/repository/redis"
	"apartment-manager-backend/pkg/hasher"
	"context"
	"errors"
	"time"
)

type LoginService struct {
	userRepo     postgres.UserInterface
	refreshRepo  redis.RefreshTokenInterFace
	tokenService jwt.TokenServiceInterface
	hasher       *hasher.BcryptHasher
	refreshTTL   time.Duration
}

func NewLoginService(
	userRepo postgres.UserInterface,
	refreshRepo redis.RefreshTokenInterFace,
	tokenService jwt.TokenServiceInterface,
	hasher *hasher.BcryptHasher,
) *LoginService {
	return &LoginService{
		userRepo:     userRepo,
		refreshRepo:  refreshRepo,
		tokenService: tokenService,
		hasher:       hasher,
		refreshTTL:   7 * 24 * time.Hour,
	}
}

func (s *LoginService) Execute(ctx context.Context, input dto.LoginInput) (*dto.LoginOutput, error) {
	user, err := s.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid_credentials")
	}

	if !s.hasher.Compare(input.Password, user.Password) {
		return nil, errors.New("invalid_credentials")
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("failed_to_generate_access_token")
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("failed_to_generate_refresh_token")
	}

	if s.refreshRepo != nil {
		err := s.refreshRepo.Save(ctx, user.ID.String(), refreshToken, s.refreshTTL)
		if err != nil {
			return nil, errors.New("failed_to_store_session")
		}
	}

	return &dto.LoginOutput{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "login_success",
	}, nil
}
