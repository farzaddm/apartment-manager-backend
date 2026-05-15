package service

import (
	"apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"errors"
)

type ProfileService struct {
	userRepo postgres.UserInterface
}

func NewProfileService(userRepo postgres.UserInterface) *ProfileService {
	return &ProfileService{userRepo: userRepo}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user_not_found")
	}

	return map[string]interface{}{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"username":   user.Username,
		"email":      user.Email,
		"phone":      user.Phone,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	}, nil
}
