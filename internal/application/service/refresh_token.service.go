package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/domain/repository/redis"
	"context"
	"errors"
	"time"
)

type RefreshTokenService struct {
	refreshRepo  redis.RefreshTokenInterFace
	tokenService jwt.TokenServiceInterface
	refreshTTL   time.Duration
}

func NewRefreshTokenService(
	refreshRepo redis.RefreshTokenInterFace,
	tokenService jwt.TokenServiceInterface,
) *RefreshTokenService {
	return &RefreshTokenService{
		refreshRepo:  refreshRepo,
		tokenService: tokenService,
		refreshTTL:   7 * 24 * time.Hour,
	}
}

func (s *RefreshTokenService) Execute(ctx context.Context, input dto.RefreshInput) (*dto.LoginOutput, error) {
	claims, err := s.tokenService.ValidateToken(input.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid_or_expired_refresh_token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid_token_type")
	}

	storedToken, err := s.refreshRepo.Get(ctx, claims.UserID.String())
	if err != nil || storedToken != input.RefreshToken {
		return nil, errors.New("session_expired_or_revoked")
	}

	newAccessToken, err := s.tokenService.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginOutput{
		AccessToken:  newAccessToken,
		RefreshToken: input.RefreshToken,
		Message:      "access_token_renewed",
	}, nil
}
