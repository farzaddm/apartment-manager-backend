package middleware

import (
	"apartment-manager-backend/internal/domain/jwt"
	"apartment-manager-backend/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService jwt.TokenServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization_header_required", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid_authorization_format", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid_or_expired_token", err)
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			response.Error(c, http.StatusUnauthorized, "invalid_token_type", nil)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("role", claims.Role)

		c.Next()
	}
}
