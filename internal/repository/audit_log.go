package repository

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AuditLogRepository) List(ctx context.Context, filters map[string]interface{}, limit int, offset int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := r.db.WithContext(ctx)

	// Apply filters
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if action, ok := filters["action"]; ok {
		query = query.Where("action = ?", action)
	}
	if resourceType, ok := filters["resource_type"]; ok {
		query = query.Where("resource_type = ?", resourceType)
	}
	if resourceID, ok := filters["resource_id"]; ok {
		query = query.Where("resource_id = ?", resourceID)
	}
	if startDate, ok := filters["start_date"]; ok {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"]; ok {
		query = query.Where("created_at <= ?", endDate)
	}

	// Get total count
	query.Model(&models.AuditLog{}).Count(&total)

	// Get paginated results
	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Preload("User").
		Find(&logs).Error

	return logs, total, err
}

func (r *AuditLogRepository) GetByID(ctx context.Context, id uint) (*models.AuditLog, error) {
	var log models.AuditLog
	err := r.db.WithContext(ctx).Preload("User").First(&log, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *AuditLogRepository) GetByResourceID(ctx context.Context, resourceID string, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.WithContext(ctx).
		Where("resource_id = ?", resourceID).
		Order("created_at DESC").
		Limit(limit).
		Preload("User").
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) DeleteOlderThan(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return r.db.WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&models.AuditLog{}).Error
}

func (r *AuditLogRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	var stats struct {
		TotalLogs   int64
		TodayLogs   int64
		UniquUsers  int64
		ActionTypes map[string]int64
	}

	today := time.Now().Format("2006-01-02")

	r.db.WithContext(ctx).Model(&models.AuditLog{}).Count(&stats.TotalLogs)
	r.db.WithContext(ctx).Where("DATE(created_at) = ?", today).Count(&stats.TodayLogs)
	r.db.WithContext(ctx).Model(&models.AuditLog{}).Distinct("user_id").Count(&stats.UniquUsers)

	return map[string]interface{}{
		"total_logs":  stats.TotalLogs,
		"today_logs":  stats.TodayLogs,
		"unique_users": stats.UniquUsers,
	}, nil
}
