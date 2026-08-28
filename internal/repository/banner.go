package repository

import (
	"context"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type BannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

// Create creates a new banner
func (r *BannerRepository) Create(ctx context.Context, banner *models.Banner) error {
	return r.db.WithContext(ctx).Create(banner).Error
}

// GetByID retrieves a banner by ID
func (r *BannerRepository) GetByID(ctx context.Context, id uint) (*models.Banner, error) {
	var banner models.Banner
	err := r.db.WithContext(ctx).First(&banner, id).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// ListUserBanners retrieves all banners for a user
func (r *BannerRepository) ListUserBanners(ctx context.Context, userID uint) ([]models.Banner, error) {
	var banners []models.Banner
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("display_order ASC, created_at DESC").
		Find(&banners).Error
	return banners, err
}

// ListActiveBannersByType retrieves active banners by type
func (r *BannerRepository) ListActiveBannersByType(ctx context.Context, userID uint, bannerType string) ([]models.Banner, error) {
	var banners []models.Banner
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND is_active = true", userID, bannerType).
		Order("display_order ASC").
		Find(&banners).Error
	return banners, err
}

// Update updates a banner
func (r *BannerRepository) Update(ctx context.Context, banner *models.Banner) error {
	return r.db.WithContext(ctx).Save(banner).Error
}

// Delete deletes a banner
func (r *BannerRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Banner{}, id).Error
}

// UpdateDisplayOrder updates the display order for multiple banners
func (r *BannerRepository) UpdateDisplayOrder(ctx context.Context, bannerID uint, order int) error {
	return r.db.WithContext(ctx).
		Model(&models.Banner{}).
		Where("id = ?", bannerID).
		Update("display_order", order).Error
}

// GetBannerByS3Key retrieves banner by S3 key
func (r *BannerRepository) GetBannerByS3Key(ctx context.Context, s3Key string) (*models.Banner, error) {
	var banner models.Banner
	err := r.db.WithContext(ctx).Where("s3_key = ?", s3Key).First(&banner).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// GetBannerCount returns total count of banners for a user
func (r *BannerRepository) GetBannerCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Banner{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
