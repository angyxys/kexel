package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	serv *service.AuditService
}

func NewAuditHandler(serv *service.AuditService) *AuditHandler {
	return &AuditHandler{
		serv: serv,
	}
}

type AuditLogResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	Description  string    `json:"description"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuditListResponse struct {
	Data       []AuditLogResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

type AuditStatsResponse struct {
	TotalLogs   int64 `json:"total_logs"`
	TodayLogs   int64 `json:"today_logs"`
	UniqueUsers int64 `json:"unique_users"`
}

// ListAuditLogs retrieves paginated audit logs with optional filters
func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Filters
	filters := make(map[string]interface{})

	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	if resourceType := c.Query("resource_type"); resourceType != "" {
		filters["resource_type"] = resourceType
	}
	if resourceID := c.Query("resource_id"); resourceID != "" {
		filters["resource_id"] = resourceID
	}

	// Date range filters
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_date"] = t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filters["end_date"] = t.AddDate(0, 0, 1) // Include full day
		}
	}

	logs, total, err := h.serv.ListLogs(ctx, filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve audit logs",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	// Convert to response format
	responses := make([]AuditLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = AuditLogResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			Username:     log.User.Username,
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			Description:  log.Description,
			IPAddress:    log.IPAddress,
			CreatedAt:    log.CreatedAt,
		}
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, AuditListResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(totalPages),
	})
}

// GetAuditStats retrieves audit log statistics
func (h *AuditHandler) GetAuditStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	stats, err := h.serv.GetStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve statistics",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ExportAuditLogs exports audit logs as CSV
func (h *AuditHandler) ExportAuditLogs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	// Filters
	filters := make(map[string]interface{})

	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	if resourceType := c.Query("resource_type"); resourceType != "" {
		filters["resource_type"] = resourceType
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_date"] = t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filters["end_date"] = t.AddDate(0, 0, 1)
		}
	}

	// Get all logs (no pagination for export)
	logs, _, err := h.serv.ListLogs(ctx, filters, 1, 10000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to export logs",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	// Build CSV
	csv := "ID,User,Action,Resource Type,Resource ID,IP Address,Timestamp\n"
	for _, log := range logs {
		csv += fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s\n",
			log.ID,
			log.User.Username,
			log.Action,
			log.ResourceType,
			log.ResourceID,
			log.IPAddress,
			log.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	// Return CSV file
	c.Header("Content-Disposition", "attachment; filename=audit-logs.csv")
	c.Header("Content-Type", "text/csv")
	c.String(http.StatusOK, csv)
}

// GetResourceHistory retrieves all logs for a specific resource
func (h *AuditHandler) GetResourceHistory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	resourceID := c.Param("id")

	logs, err := h.serv.GetLogsByResource(ctx, resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve resource history",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	responses := make([]AuditLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = AuditLogResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			Username:     log.User.Username,
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			Description:  log.Description,
			IPAddress:    log.IPAddress,
			CreatedAt:    log.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, responses)
}
