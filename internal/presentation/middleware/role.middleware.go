package middleware

import (
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RolesAuthorize(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "unauthorized_access")
			c.Abort()
			return
		}

		hasAccess := false
		for _, role := range allowedRoles {
			if role == userRole.(string) {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			response.Error(c, http.StatusForbidden, "forbidden_you_dont_have_access")
			c.Abort()
			return
		}

		c.Next()
	}
}
