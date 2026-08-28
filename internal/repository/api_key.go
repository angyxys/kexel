package repository

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// Create creates a new API key
func (r *APIKeyRepository) Create(ctx context.Context, apiKey *models.APIKey) error {
	return r.db.WithContext(ctx).Create(apiKey).Error
}

// GetByKey retrieves an API key by key value
func (r *APIKeyRepository) GetByKey(ctx context.Context, key string) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.WithContext(ctx).Where("key = ? AND is_active = true", key).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// GetByID retrieves an API key by ID
func (r *APIKeyRepository) GetByID(ctx context.Context, id uint) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.WithContext(ctx).First(&apiKey, id).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// ListUserAPIKeys retrieves all API keys for a user
func (r *APIKeyRepository) ListUserAPIKeys(ctx context.Context, userID uint) ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&apiKeys).Error
	return apiKeys, err
}

// ListActiveUserAPIKeys retrieves active API keys for a user
func (r *APIKeyRepository) ListActiveUserAPIKeys(ctx context.Context, userID uint) ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = true", userID).
		Order("created_at DESC").
		Find(&apiKeys).Error
	return apiKeys, err
}

// UpdateLastUsed updates the last used timestamp
func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, keyID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("id = ?", keyID).
		Update("last_used_at", time.Now()).Error
}

// Revoke marks an API key as inactive
func (r *APIKeyRepository) Revoke(ctx context.Context, keyID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("id = ?", keyID).
		Update("is_active", false).Error
}

// Delete deletes an API key
func (r *APIKeyRepository) Delete(ctx context.Context, keyID uint) error {
	return r.db.WithContext(ctx).Delete(&models.APIKey{}, keyID).Error
}

// UpdateRateLimit updates the rate limit for an API key
func (r *APIKeyRepository) UpdateRateLimit(ctx context.Context, keyID uint, rateLimit int) error {
	return r.db.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("id = ?", keyID).
		Update("rate_limit", rateLimit).Error
}

// DeleteExpiredKeys deletes API keys that have expired
func (r *APIKeyRepository) DeleteExpiredKeys(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.APIKey{}).Error
}

// GetKeyCount returns the count of active API keys for a user
func (r *APIKeyRepository) GetKeyCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("user_id = ? AND is_active = true", userID).
		Count(&count).Error
	return count, err
}

// UpdateScopes updates the scopes for an API key
func (r *APIKeyRepository) UpdateScopes(ctx context.Context, keyID uint, scopes string) error {
	return r.db.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("id = ?", keyID).
		Update("scopes", scopes).Error
}
