package service

import (
	"apartment-manager-backend/internal/domain/repository/redis"
	"context"
	"errors"
)

type LogoutService struct {
	refreshRepo redis.RefreshTokenInterFace
}

func NewLogoutService(refreshRepo redis.RefreshTokenInterFace) *LogoutService {
	return &LogoutService{
		refreshRepo: refreshRepo,
	}
}

func (s *LogoutService) Execute(ctx context.Context, userID string) error {
	_, err := s.refreshRepo.Get(ctx, userID)
	if err != nil {
		return errors.New("no_active_session_found")
	}

	err = s.refreshRepo.Delete(ctx, userID)
	if err != nil {
		return errors.New("failed_to_clear_session")
	}

	return nil
}
