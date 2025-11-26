package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	ValidateToken() gin.HandlerFunc
}

func SetupRoutes(r *gin.Engine, auth AuthMiddleware) {
	// Public Routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/api/directory/lookup", GetOrganizationByDomain)

	// Protected Routes
	protected := r.Group("/api")
	protected.Use(auth.ValidateToken())
	{
		protected.GET("/messages", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "This is a protected message from the Go Backend.",
			})
		})
	}
}
