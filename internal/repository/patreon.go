package repository

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type PatreonIntegrationRepository struct {
	db *gorm.DB
}

func NewPatreonIntegrationRepository(db *gorm.DB) *PatreonIntegrationRepository {
	return &PatreonIntegrationRepository{db: db}
}

// Create creates a new Patreon integration
func (r *PatreonIntegrationRepository) Create(ctx context.Context, integration *models.PatreonIntegration) error {
	return r.db.WithContext(ctx).Create(integration).Error
}

// GetByUserID retrieves Patreon integration by user ID
func (r *PatreonIntegrationRepository) GetByUserID(ctx context.Context, userID uint) (*models.PatreonIntegration, error) {
	var integration models.PatreonIntegration
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&integration).Error
	if err != nil {
		return nil, err
	}
	return &integration, nil
}

// Update updates a Patreon integration
func (r *PatreonIntegrationRepository) Update(ctx context.Context, integration *models.PatreonIntegration) error {
	return r.db.WithContext(ctx).Save(integration).Error
}

// Delete deletes a Patreon integration
func (r *PatreonIntegrationRepository) Delete(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.PatreonIntegration{}).Error
}

// UpdateLastSync updates the last sync timestamp
func (r *PatreonIntegrationRepository) UpdateLastSync(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.PatreonIntegration{}).
		Where("user_id = ?", userID).
		Update("last_sync_at", time.Now()).Error
}

// ListEnabledIntegrations retrieves all enabled Patreon integrations
func (r *PatreonIntegrationRepository) ListEnabledIntegrations(ctx context.Context) ([]models.PatreonIntegration, error) {
	var integrations []models.PatreonIntegration
	err := r.db.WithContext(ctx).
		Where("is_enabled = true").
		Find(&integrations).Error
	return integrations, err
}

// RateLimitRuleRepository handles rate limit rules
type RateLimitRuleRepository struct {
	db *gorm.DB
}

func NewRateLimitRuleRepository(db *gorm.DB) *RateLimitRuleRepository {
	return &RateLimitRuleRepository{db: db}
}

// Create creates a new rate limit rule
func (r *RateLimitRuleRepository) Create(ctx context.Context, rule *models.RateLimitRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

// GetByIP retrieves rate limit rule by IP
func (r *RateLimitRuleRepository) GetByIP(ctx context.Context, ip string) (*models.RateLimitRule, error) {
	var rule models.RateLimitRule
	err := r.db.WithContext(ctx).
		Where("ip_address = ? AND is_active = true", ip).
		First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Update updates a rate limit rule
func (r *RateLimitRuleRepository) Update(ctx context.Context, rule *models.RateLimitRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

// Delete deletes a rate limit rule
func (r *RateLimitRuleRepository) Delete(ctx context.Context, ip string) error {
	return r.db.WithContext(ctx).Where("ip_address = ?", ip).Delete(&models.RateLimitRule{}).Error
}

// RateLimitLogRepository handles rate limit logs
type RateLimitLogRepository struct {
	db *gorm.DB
}

func NewRateLimitLogRepository(db *gorm.DB) *RateLimitLogRepository {
	return &RateLimitLogRepository{db: db}
}

// Create creates a new rate limit log
func (r *RateLimitLogRepository) Create(ctx context.Context, log *models.RateLimitLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetRecentBlocks retrieves recent rate limit blocks
func (r *RateLimitLogRepository) GetRecentBlocks(ctx context.Context, limit int) ([]models.RateLimitLog, error) {
	var logs []models.RateLimitLog
	err := r.db.WithContext(ctx).
		Order("blocked_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetBlockCountForIP gets count of blocks for an IP in the last hour
func (r *RateLimitLogRepository) GetBlockCountForIP(ctx context.Context, ip string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.RateLimitLog{}).
		Where("ip_address = ? AND blocked_at > ?", ip, time.Now().Add(-time.Hour)).
		Count(&count).Error
	return count, err
}
