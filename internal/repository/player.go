package repository

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type PlayerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

func (r *PlayerRepository) GetPlayerByID(ctx context.Context, playerID string) (*models.Player, error) {
	var player models.Player
	err := r.db.WithContext(ctx).First(&player, "vrchat_id = ?", playerID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *PlayerRepository) AddPlayer(ctx context.Context, player *models.Player) error {
	return r.db.WithContext(ctx).Create(player).Error
}

func (r *PlayerRepository) UpdatePlayer(ctx context.Context, player *models.Player) error {
	return r.db.WithContext(ctx).Save(player).Error
}

func (r *PlayerRepository) DeletePlayer(ctx context.Context, playerID string) error {
	return r.db.WithContext(ctx).Delete(&models.Player{}, "vrchat_id = ?", playerID).Error
}

func (r *PlayerRepository) ListPlayers(ctx context.Context) ([]models.Player, error) {
	var players []models.Player
	err := r.db.WithContext(ctx).Find(&players).Error
	return players, err
}

// SearchPlayers performs full-text search on players
func (r *PlayerRepository) SearchPlayers(ctx context.Context, query string, limit int, offset int) ([]models.Player, int64, error) {
	var players []models.Player
	var total int64

	db := r.db.WithContext(ctx)

	// Search by VRChat ID (exact or partial match)
	if query != "" {
		db = db.Where("vrchat_id ILIKE ?", "%"+query+"%")
	}

	// Get total count
	db.Model(&models.Player{}).Count(&total)

	// Get paginated results
	err := db.
		Order("vrchat_id ASC").
		Limit(limit).
		Offset(offset).
		Find(&players).Error

	return players, total, err
}

// FilterPlayers applies multiple filters to players
type PlayerFilters struct {
	Search    string
	Roles     []string // Filter by specific roles
	IsBanned  *bool    // Filter by ban status
	SortBy    string   // vrchat_id, created_at, updated_at
	SortOrder string   // ASC, DESC
}

func (r *PlayerRepository) FilterPlayers(ctx context.Context, filters PlayerFilters, limit int, offset int) ([]models.Player, int64, error) {
	var players []models.Player
	var total int64

	db := r.db.WithContext(ctx)

	// Search filter
	if filters.Search != "" {
		db = db.Where("vrchat_id ILIKE ?", "%"+filters.Search+"%")
	}

	// Ban status filter
	if filters.IsBanned != nil {
		db = db.Where("is_banned = ?", *filters.IsBanned)
	}

	// Roles filter (checks if any role matches)
	if len(filters.Roles) > 0 {
		db = db.Where("role && ?", filters.Roles) // PostgreSQL array overlap operator
	}

	// Get total count
	db.Model(&models.Player{}).Count(&total)

	// Sorting
	sortBy := "vrchat_id"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "ASC"
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}

	// Get paginated results
	err := db.
		Order(sortBy + " " + sortOrder).
		Limit(limit).
		Offset(offset).
		Find(&players).Error

	return players, total, err
}

// GetBannedPlayers returns all currently banned players
func (r *PlayerRepository) GetBannedPlayers(ctx context.Context) ([]models.Player, error) {
	var players []models.Player
	err := r.db.WithContext(ctx).
		Where("is_banned = true").
		Order("ban_expires_at ASC NULLS LAST").
		Find(&players).Error
	return players, err
}

// GetExpiringSoonBans returns players whose bans expire within the specified duration
func (r *PlayerRepository) GetExpiringSoonBans(ctx context.Context, duration time.Duration) ([]models.Player, error) {
	var players []models.Player
	now := time.Now()
	expiryThreshold := now.Add(duration)

	err := r.db.WithContext(ctx).
		Where("is_banned = true AND ban_expires_at IS NOT NULL").
		Where("ban_expires_at > ? AND ban_expires_at <= ?", now, expiryThreshold).
		Order("ban_expires_at ASC").
		Find(&players).Error
	return players, err
}

// UnbanExpiredPlayers unbans all players whose ban has expired
func (r *PlayerRepository) UnbanExpiredPlayers(ctx context.Context) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.Player{}).
		Where("is_banned = true AND ban_expires_at IS NOT NULL AND ban_expires_at <= ?", now).
		Updates(map[string]interface{}{
			"is_banned":       false,
			"ban_reason":      "",
			"ban_expires_at":  nil,
			"banned_at":       nil,
		})
	return result.RowsAffected, result.Error
}

func (r *PlayerRepository) BanUnbanPlayer(ctx context.Context, playerID string, ban bool) error {
	return r.db.WithContext(ctx).Model(&models.Player{}).Where("vrchat_id = ?", playerID).Update("is_banned", ban).Error
}

func (r *PlayerRepository) AllVip(ctx context.Context) []models.Player {
	var players []models.Player
	vip := []models.Role{
		models.ROLE_VIP,
	}
	err := r.db.WithContext(ctx).Where("role IN ?", vip).Find(&players).Error
	if err != nil {
		return make([]models.Player, 0)
	}
	return players
}

func (r *PlayerRepository) AllBanned(ctx context.Context) []models.Player {
	var players []models.Player
	err := r.db.WithContext(ctx).Where("is_banned = true").Find(&players).Error
	if err != nil {
		return make([]models.Player, 0)
	}
	return players
}

func (r *PlayerRepository) AllModerator(ctx context.Context) []models.Player {
	var players []models.Player
	vip := []models.Role{
		models.ROLE_MOD,
	}
	err := r.db.WithContext(ctx).Where("role IN ?", vip).Find(&players).Error
	if err != nil {
		return make([]models.Player, 0)
	}
	return players
}

func (r *PlayerRepository) AllOwner(ctx context.Context) []models.Player {
	var players []models.Player
	vip := []models.Role{
		models.ROLE_OWNER,
	}
	err := r.db.WithContext(ctx).Where("role IN ?", vip).Find(&players).Error
	if err != nil {
		return make([]models.Player, 0)
	}
	return players
}
