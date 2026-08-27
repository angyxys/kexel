package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type AuditService struct {
	repo *repository.AuditLogRepository
}

func NewAuditService(repo *repository.AuditLogRepository) *AuditService {
	return &AuditService{repo: repo}
}

// LogAction logs an action
func (s *AuditService) LogAction(ctx context.Context, log *models.AuditLog) error {
	return s.repo.Create(ctx, log)
}

// ListLogs retrieves audit logs with pagination and filters
func (s *AuditService) ListLogs(ctx context.Context, filters map[string]interface{}, page int, pageSize int) ([]models.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	return s.repo.List(ctx, filters, pageSize, offset)
}

// GetLogsByResource retrieves all logs for a specific resource
func (s *AuditService) GetLogsByResource(ctx context.Context, resourceID string) ([]models.AuditLog, error) {
	return s.repo.GetByResourceID(ctx, resourceID, 50)
}

// GetStats retrieves audit log statistics
func (s *AuditService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.GetStats(ctx)
}

// ArchiveOldLogs deletes logs older than specified days (cleanup)
func (s *AuditService) ArchiveOldLogs(ctx context.Context, days int) error {
	if days < 7 {
		days = 7 // Minimum 7 days
	}
	return s.repo.DeleteOlderThan(ctx, days)
}

// Helper to create audit log entry
func NewAuditLog(userID uint, action string, resourceType string, resourceID string, ipAddress string, userAgent string) *models.AuditLog {
	return &models.AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    time.Now(),
	}
}

// Helper functions for JSON serialization
func SetOldValue(al *models.AuditLog, value interface{}) error {
	if value == nil {
		al.OldValue = ""
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	al.OldValue = string(data)
	return nil
}

func SetNewValue(al *models.AuditLog, value interface{}) error {
	if value == nil {
		al.NewValue = ""
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	al.NewValue = string(data)
	return nil
}

func GetOldValue(al *models.AuditLog) (map[string]interface{}, error) {
	var result map[string]interface{}
	if al.OldValue == "" {
		return result, nil
	}
	err := json.Unmarshal([]byte(al.OldValue), &result)
	return result, err
}

func GetNewValue(al *models.AuditLog) (map[string]interface{}, error) {
	var result map[string]interface{}
	if al.NewValue == "" {
		return result, nil
	}
	err := json.Unmarshal([]byte(al.NewValue), &result)
	return result, err
}
