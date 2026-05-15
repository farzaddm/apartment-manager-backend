package service

import (
	"apartment-manager-backend/internal/domain/entity"
	domainJwt "apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"apartment-manager-backend/internal/domain/repository/redis"
	"context"
	"errors"
	"time"
)

type VerifyOTPOutput struct {
	User         *entity.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	Message      string       `json:"message"`
}

type VerifyOTPService struct {
	otpRepo      redis.OTPInterFace
	userRepo     postgres.UserInterface
	refreshRepo  redis.RefreshTokenInterFace
	tokenService domainJwt.TokenServiceInterface
	refreshTTL   time.Duration
}

func NewVerifyOTPService(
	otpRepo redis.OTPInterFace,
	userRepo postgres.UserInterface,
	refreshRepo redis.RefreshTokenInterFace,
	tokenService domainJwt.TokenServiceInterface,
) *VerifyOTPService {
	return &VerifyOTPService{
		otpRepo:      otpRepo,
		userRepo:     userRepo,
		refreshRepo:  refreshRepo,
		tokenService: tokenService,
		refreshTTL:   7 * 24 * time.Hour,
	}
}

func (s *VerifyOTPService) Execute(ctx context.Context, phone, code string) (*VerifyOTPOutput, error) {
	storedOTP, err := s.otpRepo.Get(phone)
	if err != nil {
		return nil, errors.New("otp_not_found_or_expired")
	}

	if storedOTP != code {
		return nil, errors.New("invalid_otp")
	}

	_ = s.otpRepo.Delete(phone)

	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if err := s.otpRepo.SetVerified(phone, 5*time.Minute); err != nil {
			return nil, errors.New("failed_to_store_verification_status")
		}

		return &VerifyOTPOutput{
			User:    nil,
			Message: "account_not_found",
		}, nil
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	if s.refreshRepo != nil {
		if err := s.refreshRepo.Save(ctx, user.ID.String(), refreshToken, s.refreshTTL); err != nil {
			return nil, err
		}
	}

	return &VerifyOTPOutput{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "login_success",
	}, nil
}
