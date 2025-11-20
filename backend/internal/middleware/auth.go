package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	Domain   string
	Audience string
}

func NewAuthMiddleware(domain, audience string) *AuthMiddleware {
	return &AuthMiddleware{
		Domain:   domain,
		Audience: audience,
	}
}

func (m *AuthMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		// token := parts[1]

		// TODO: Implement actual JWT validation using a library like github.com/auth0/go-jwt-middleware/v2
		// For MVP structure, we are mocking the validation pass to ensure connectivity first.
		// In a real implementation, we would fetch JWKS from Auth0 and validate the signature.

		// MOCK VALIDATION FOR NOW to allow testing without full Auth0 setup immediately
		// if token == "invalid" {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// 	return
		// }

		c.Next()
	}
}
