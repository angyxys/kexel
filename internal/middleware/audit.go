package middleware

import (
	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

// AuditMiddleware logs all mutations to the audit log
func AuditMiddleware(auditServ *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log mutations (POST, PUT, DELETE)
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			userID, exists := c.Get("user_id")
			if !exists {
				userID = uint(0) // Anonymous if not logged in
			}

			// Create audit log entry
			log := &models.AuditLog{
				UserID:    userID.(uint),
				Action:    c.Request.Method,
				IPAddress: c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
			}

			// Parse resource from URL path
			// Example: /web/player/:id -> resourceType=player, resourceID=:id
			if c.Request.URL.Path != "" {
				// This will be populated after the route is matched
				c.Set("audit_log", log)
			}
		}

		// Process request
		c.Next()

		// After processing, log if it was a mutation
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			if auditLog, exists := c.Get("audit_log"); exists {
				log := auditLog.(*models.AuditLog)

				// Add response status
				log.Description = c.GetString("audit_description")
				if log.Description == "" {
					log.Description = "API " + c.Request.Method + " " + c.Request.URL.Path
				}

				// Only log if status indicates action was processed
				if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
					auditServ.LogAction(c.Request.Context(), log)
				}
			}
		}
	}
}

// SetAuditInfo helper to set audit log information in context
func SetAuditInfo(c *gin.Context, resourceType string, resourceID string, description string) {
	if auditLog, exists := c.Get("audit_log"); exists {
		log := auditLog.(*models.AuditLog)
		log.ResourceType = resourceType
		log.ResourceID = resourceID
		log.Description = description
		c.Set("audit_log", log)
	}
	c.Set("audit_description", description)
}
