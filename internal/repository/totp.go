package repository

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type TOTPRepository struct {
	db *gorm.DB
}

func NewTOTPRepository(db *gorm.DB) *TOTPRepository {
	return &TOTPRepository{db: db}
}

// Create creates a new TOTP secret
func (r *TOTPRepository) Create(ctx context.Context, totpSecret *models.TOTPSecret) error {
	// Delete any existing TOTP secret for this user first
	r.db.WithContext(ctx).Where("user_id = ?", totpSecret.UserID).Delete(&models.TOTPSecret{})

	return r.db.WithContext(ctx).Create(totpSecret).Error
}

// GetByUserID retrieves TOTP secret by user ID
func (r *TOTPRepository) GetByUserID(ctx context.Context, userID uint) (*models.TOTPSecret, error) {
	var totpSecret models.TOTPSecret
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&totpSecret).Error
	if err != nil {
		return nil, err
	}
	return &totpSecret, nil
}

// EnableTOTP enables TOTP for a user
func (r *TOTPRepository) EnableTOTP(ctx context.Context, userID uint, enabledAt *time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.TOTPSecret{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"is_enabled": true,
			"enabled_at": enabledAt,
		}).Error
}

// DisableTOTP disables TOTP for a user
func (r *TOTPRepository) DisableTOTP(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.TOTPSecret{}).
		Where("user_id = ?", userID).
		Update("is_enabled", false).Error
}

// UpdateBackupCodes updates the backup codes for a user
func (r *TOTPRepository) UpdateBackupCodes(ctx context.Context, userID uint, backupCodes string) error {
	return r.db.WithContext(ctx).
		Model(&models.TOTPSecret{}).
		Where("user_id = ?", userID).
		Update("backup_codes", backupCodes).Error
}

// Delete deletes TOTP secret for a user
func (r *TOTPRepository) Delete(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.TOTPSecret{}).Error
}

// UpdateLastUsed updates the last used timestamp
func (r *TOTPRepository) UpdateLastUsed(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.TOTPSecret{}).
		Where("user_id = ?", userID).
		Update("last_used_at", time.Now()).Error
}
