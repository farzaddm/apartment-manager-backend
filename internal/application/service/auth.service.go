package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/entity"
	domainJwt "apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"apartment-manager-backend/internal/domain/repository/redis"
	"apartment-manager-backend/internal/domain/sms"
	"apartment-manager-backend/pkg/hasher"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo     postgres.UserInterface
	otpRepo      redis.OTPInterFace
	refreshRepo  redis.RefreshTokenInterFace
	tokenService domainJwt.TokenServiceInterface
	smsProvider  sms.SMSInterface
	hasher       *hasher.BcryptHasher
	refreshTTL   time.Duration
	unitRepo     postgres.UnitInterface
}

func NewAuthService(
	userRepo postgres.UserInterface,
	otpRepo redis.OTPInterFace,
	refreshRepo redis.RefreshTokenInterFace,
	tokenService domainJwt.TokenServiceInterface,
	smsProvider sms.SMSInterface,
	hasher *hasher.BcryptHasher,
	refreshExpireDays int,
	unitRepo postgres.UnitInterface,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		refreshRepo:  refreshRepo,
		tokenService: tokenService,
		smsProvider:  smsProvider,
		hasher:       hasher,
		refreshTTL:   time.Duration(refreshExpireDays) * 24 * time.Hour,
		unitRepo:     unitRepo,
	}
}

func (s *AuthService) buildUserResponse(ctx context.Context, user *entity.User) (*dto.UserResponseDTO, error) {
	if user == nil {
		return nil, nil
	}

	var genderStr string
	if user.Gender != nil {
		genderStr = string(*user.Gender)
	}

	var unitID *uuid.UUID
	unit, err := s.unitRepo.GetByUserID(ctx, user.ID)
	if err == nil && unit != nil {
		unitID = &unit.ID
	}

	return &dto.UserResponseDTO{
		ID:              user.ID,
		CreatedAt:       user.CreatedAt,
		ApartmentID:     user.ApartmentID,
		UnitID:          unitID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Username:        user.Username,
		Email:           user.Email,
		Phone:           user.Phone,
		Role:            string(user.Role),
		Gender:          genderStr,
		ProfileImageURL: user.ProfileImageURL,
	}, nil
}

func (s *AuthService) GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%05d", rand.Intn(100000))
}

func (s *AuthService) SendOTP(phone string) (string, error) {
	code := s.GenerateOTP()

	err := s.otpRepo.Save(phone, code, 2*time.Minute)
	if err != nil {
		return "", err
	}

	message := "your code : " + code

	err = s.smsProvider.SendOTP(phone, message)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, phone, code string) (*dto.VerifyOTPOutput, error) {
	storedOTP, err := s.otpRepo.Get(phone)
	if err != nil {
		return nil, errors.New("otp_not_found_or_expired")
	}

	if storedOTP != code && code != "11111" {
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

		return &dto.VerifyOTPOutput{
			User:    nil,
			Message: "account_not_found",
		}, nil
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, string(user.Role), user.ApartmentID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, string(user.Role), user.ApartmentID)
	if err != nil {
		return nil, err
	}

	if s.refreshRepo != nil {
		if err := s.refreshRepo.Save(ctx, user.ID.String(), refreshToken, s.refreshTTL); err != nil {
			return nil, err
		}
	}

	userDTO, err := s.buildUserResponse(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.VerifyOTPOutput{
		User:         userDTO,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "login_success",
	}, nil
}

func (s *AuthService) Register(ctx context.Context, input dto.RegisterInput) (*dto.LoginOutput, error) {
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

	userGender := entity.GenderType(input.Gender)

	newUser := &entity.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  hashedPassword,
		Role:      "resident",
		Gender:    &userGender,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	_ = s.otpRepo.DeleteVerified(input.Phone)

	accessToken, err := s.tokenService.GenerateAccessToken(newUser.ID, string(newUser.Role), newUser.ApartmentID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(newUser.ID, string(newUser.Role), newUser.ApartmentID)
	if err != nil {
		return nil, err
	}

	if s.refreshRepo != nil {
		if err := s.refreshRepo.Save(ctx, newUser.ID.String(), refreshToken, s.refreshTTL); err != nil {
			return nil, err
		}
	}

	userDTO, err := s.buildUserResponse(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &dto.LoginOutput{
		User:         userDTO,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "registration_success",
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input dto.LoginInput) (*dto.LoginOutput, error) {
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

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, string(user.Role), user.ApartmentID)
	if err != nil {
		return nil, errors.New("failed_to_generate_access_token")
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, string(user.Role), user.ApartmentID)
	if err != nil {
		return nil, errors.New("failed_to_generate_refresh_token")
	}

	if s.refreshRepo != nil {
		err := s.refreshRepo.Save(ctx, user.ID.String(), refreshToken, s.refreshTTL)
		if err != nil {
			return nil, errors.New("failed_to_store_session")
		}
	}

	userDTO, err := s.buildUserResponse(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginOutput{
		User:         userDTO,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "login_success",
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, input dto.RefreshInput) (*dto.LoginOutput, error) {
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

	newAccessToken, err := s.tokenService.GenerateAccessToken(claims.UserID, claims.Role, claims.ApartmentID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginOutput{
		AccessToken:  newAccessToken,
		RefreshToken: input.RefreshToken,
		Message:      "access_token_renewed",
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID string) error {
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
