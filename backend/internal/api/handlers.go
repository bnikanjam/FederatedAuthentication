package api

import (
	"federation-auth/internal/db"
	"federation-auth/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrganizationByDomain(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain parameter is required"})
		return
	}

	var org models.Organization
	result := db.DB.Where("domain = ?", domain).First(&org)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organization_id": org.Auth0OrgID,
		"display_name":    org.DisplayName,
	})
}
