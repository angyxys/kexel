package repository

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type InvitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) *InvitationRepository {
	return &InvitationRepository{db: db}
}

// Create creates a new invitation code
func (r *InvitationRepository) Create(ctx context.Context, invitation *models.InvitationCode) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

// GetByCode retrieves an invitation by code
func (r *InvitationRepository) GetByCode(ctx context.Context, code string) (*models.InvitationCode, error) {
	var invitation models.InvitationCode
	err := r.db.WithContext(ctx).Where("code = ?", code).Preload("Creator").First(&invitation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// GetByID retrieves an invitation by ID
func (r *InvitationRepository) GetByID(ctx context.Context, id uint) (*models.InvitationCode, error) {
	var invitation models.InvitationCode
	err := r.db.WithContext(ctx).Preload("Creator").First(&invitation, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// ListByCreator retrieves all invitations created by a user
func (r *InvitationRepository) ListByCreator(ctx context.Context, createdBy uint) ([]models.InvitationCode, error) {
	var invitations []models.InvitationCode
	err := r.db.WithContext(ctx).
		Where("created_by = ?", createdBy).
		Order("created_at DESC").
		Find(&invitations).Error
	return invitations, err
}

// ListActive retrieves all active invitations
func (r *InvitationRepository) ListActive(ctx context.Context) ([]models.InvitationCode, error) {
	var invitations []models.InvitationCode
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("is_active = true AND (expires_at IS NULL OR expires_at > ?)", now).
		Where("(max_uses = -1 OR uses < max_uses)").
		Order("created_at DESC").
		Preload("Creator").
		Find(&invitations).Error
	return invitations, err
}

// IncrementUses increments the use count for an invitation
func (r *InvitationRepository) IncrementUses(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.InvitationCode{}).
		Where("id = ?", id).
		Update("uses", gorm.Expr("uses + 1")).Error
}

// Revoke deactivates an invitation code
func (r *InvitationRepository) Revoke(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.InvitationCode{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// DeleteExpired deletes expired invitations
func (r *InvitationRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at <= ?", now).
		Where("is_active = false").
		Delete(&models.InvitationCode{}).Error
}

// IsValid checks if an invitation code is valid and can be used
func (r *InvitationRepository) IsValid(ctx context.Context, code string) (bool, error) {
	var count int64
	now := time.Now()

	err := r.db.WithContext(ctx).
		Model(&models.InvitationCode{}).
		Where("code = ?", code).
		Where("is_active = true").
		Where("(expires_at IS NULL OR expires_at > ?)", now).
		Where("(max_uses = -1 OR uses < max_uses)").
		Count(&count).Error

	return count > 0, err
}
