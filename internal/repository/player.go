package repository

import (
	"context"
	"log"

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

func (r *PlayerRepository) GetPlayerByID(ctx context.Context, playerID string) *models.Player {
	var player models.Player
	if err := r.db.WithContext(ctx).First(&player, "vrchat_id = ?", playerID).Error; err != nil {
		log.Printf("error when find player: %v with id: %s", err, playerID)
		return nil
	}
	return &player
}

func (r *PlayerRepository) AddPlayer(ctx context.Context, player *models.Player) {
	r.db.WithContext(ctx).Save(&player)
}

func (r *PlayerRepository) BanUnbanPlayer(ctx context.Context, playerID string, ban bool) {
	r.db.WithContext(ctx).Model(&models.Player{}).Where("vrchat_id = ?", playerID).Update("is_banned", ban)
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
