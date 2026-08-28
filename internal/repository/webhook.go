package repository

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

// Create creates a new webhook
func (r *WebhookRepository) Create(ctx context.Context, webhook *models.Webhook) error {
	return r.db.WithContext(ctx).Create(webhook).Error
}

// GetByID retrieves a webhook by ID
func (r *WebhookRepository) GetByID(ctx context.Context, id uint) (*models.Webhook, error) {
	var webhook models.Webhook
	err := r.db.WithContext(ctx).First(&webhook, id).Error
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// ListUserWebhooks retrieves all webhooks for a user
func (r *WebhookRepository) ListUserWebhooks(ctx context.Context, userID uint) ([]models.Webhook, error) {
	var webhooks []models.Webhook
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&webhooks).Error
	return webhooks, err
}

// ListActiveWebhooks retrieves all active webhooks
func (r *WebhookRepository) ListActiveWebhooks(ctx context.Context) ([]models.Webhook, error) {
	var webhooks []models.Webhook
	err := r.db.WithContext(ctx).
		Where("is_active = true").
		Find(&webhooks).Error
	return webhooks, err
}

// Update updates a webhook
func (r *WebhookRepository) Update(ctx context.Context, webhook *models.Webhook) error {
	return r.db.WithContext(ctx).Save(webhook).Error
}

// Delete deletes a webhook
func (r *WebhookRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Webhook{}, id).Error
}

// UpdateFailureCount updates the failure count
func (r *WebhookRepository) UpdateFailureCount(ctx context.Context, id uint, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.Webhook{}).
		Where("id = ?", id).
		Update("failure_count", count).Error
}

// UpdateLastTried updates the last tried timestamp
func (r *WebhookRepository) UpdateLastTried(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Webhook{}).
		Where("id = ?", id).
		Update("last_tried_at", time.Now()).Error
}

// UpdateLastSuccess updates the last success timestamp
func (r *WebhookRepository) UpdateLastSuccess(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Webhook{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_success_at": time.Now(),
			"failure_count":   0,
		}).Error
}

// WebhookEventRepository handles webhook event operations
type WebhookEventRepository struct {
	db *gorm.DB
}

func NewWebhookEventRepository(db *gorm.DB) *WebhookEventRepository {
	return &WebhookEventRepository{db: db}
}

// Create creates a new webhook event
func (r *WebhookEventRepository) Create(ctx context.Context, event *models.WebhookEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetByID retrieves a webhook event by ID
func (r *WebhookEventRepository) GetByID(ctx context.Context, id uint) (*models.WebhookEvent, error) {
	var event models.WebhookEvent
	err := r.db.WithContext(ctx).First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// ListWebhookEvents retrieves all events for a webhook
func (r *WebhookEventRepository) ListWebhookEvents(ctx context.Context, webhookID uint, limit int, offset int) ([]models.WebhookEvent, error) {
	var events []models.WebhookEvent
	err := r.db.WithContext(ctx).
		Where("webhook_id = ?", webhookID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

// ListPendingEvents retrieves events pending retry
func (r *WebhookEventRepository) ListPendingEvents(ctx context.Context) ([]models.WebhookEvent, error) {
	var events []models.WebhookEvent
	err := r.db.WithContext(ctx).
		Where("is_delivered = false AND next_retry <= ?", time.Now()).
		Order("created_at ASC").
		Find(&events).Error
	return events, err
}

// UpdateEvent updates a webhook event
func (r *WebhookEventRepository) UpdateEvent(ctx context.Context, event *models.WebhookEvent) error {
	return r.db.WithContext(ctx).Save(event).Error
}

// MarkAsDelivered marks an event as delivered
func (r *WebhookEventRepository) MarkAsDelivered(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.WebhookEvent{}).
		Where("id = ?", id).
		Update("is_delivered", true).Error
}

// GetEventCount returns the count of events for a webhook
func (r *WebhookEventRepository) GetEventCount(ctx context.Context, webhookID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.WebhookEvent{}).
		Where("webhook_id = ?", webhookID).
		Count(&count).Error
	return count, err
}

// GetDeliveredCount returns the count of delivered events for a webhook
func (r *WebhookEventRepository) GetDeliveredCount(ctx context.Context, webhookID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.WebhookEvent{}).
		Where("webhook_id = ? AND is_delivered = true", webhookID).
		Count(&count).Error
	return count, err
}
