package jwt

import (
	"errors"
	"time"

	domainJwt "apartment-manager-backend/internal/domain/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenService struct {
	accessSecret  string
	refreshSecret string
}

func NewTokenService(accessSecret, refreshSecret string) domainJwt.TokenServiceInterface {
	return &tokenService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (s *tokenService) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	return s.generateToken(userID, role, s.accessSecret, "access", 15*time.Minute)
}

func (s *tokenService) GenerateRefreshToken(userID uuid.UUID, role string) (string, error) {
	return s.generateToken(userID, role, s.refreshSecret, "refresh", 7*24*time.Hour)
}

func (s *tokenService) generateToken(userID uuid.UUID, role string, secret string, tokenType string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID.String(),
		"role":       role,
		"token_type": tokenType,
		"exp":        time.Now().Add(duration).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *tokenService) ValidateToken(tokenString string) (*domainJwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		tokenType, _ := claims["token_type"].(string)

		if tokenType == "refresh" {
			return []byte(s.refreshSecret), nil
		}

		return []byte(s.accessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claimsMap, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, _ := claimsMap["user_id"].(string)
		userID, _ := uuid.Parse(userIDStr)

		return &domainJwt.MapClaims{
			UserID:    userID,
			Role:      claimsMap["role"].(string),
			TokenType: claimsMap["token_type"].(string),
			ExpiresAt: int64(claimsMap["exp"].(float64)),
		}, nil
	}

	return nil, errors.New("invalid token")
}
