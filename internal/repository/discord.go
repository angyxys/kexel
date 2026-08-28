package repository

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type DiscordIntegrationRepository struct {
	db *gorm.DB
}

func NewDiscordIntegrationRepository(db *gorm.DB) *DiscordIntegrationRepository {
	return &DiscordIntegrationRepository{db: db}
}

// Create creates a new Discord integration
func (r *DiscordIntegrationRepository) Create(ctx context.Context, integration *models.DiscordIntegration) error {
	return r.db.WithContext(ctx).Create(integration).Error
}

// GetByUserID retrieves Discord integration by user ID
func (r *DiscordIntegrationRepository) GetByUserID(ctx context.Context, userID uint) (*models.DiscordIntegration, error) {
	var integration models.DiscordIntegration
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&integration).Error
	if err != nil {
		return nil, err
	}
	return &integration, nil
}

// Update updates a Discord integration
func (r *DiscordIntegrationRepository) Update(ctx context.Context, integration *models.DiscordIntegration) error {
	return r.db.WithContext(ctx).Save(integration).Error
}

// Delete deletes a Discord integration
func (r *DiscordIntegrationRepository) Delete(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.DiscordIntegration{}).Error
}

// UpdateLastSync updates the last sync timestamp
func (r *DiscordIntegrationRepository) UpdateLastSync(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.DiscordIntegration{}).
		Where("user_id = ?", userID).
		Update("last_sync_at", time.Now()).Error
}

// UpdateConnectionStatus updates the connection status
func (r *DiscordIntegrationRepository) UpdateConnectionStatus(ctx context.Context, userID uint, connected bool) error {
	return r.db.WithContext(ctx).
		Model(&models.DiscordIntegration{}).
		Where("user_id = ?", userID).
		Update("is_connected", connected).Error
}
