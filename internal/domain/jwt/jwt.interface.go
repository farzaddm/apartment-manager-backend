package jwt

import (
	"github.com/google/uuid"
)

type MapClaims struct {
	UserID    uuid.UUID
	Role      string
	TokenType string
	ExpiresAt int64
}

type TokenServiceInterface interface {
	GenerateAccessToken(userID uuid.UUID, role string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, role string) (string, error)
	ValidateToken(tokenString string) (*MapClaims, error)
}
