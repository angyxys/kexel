package service

import (
	"context"
	"fmt"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type BanService struct {
	repo *repository.PlayerRepository
}

func NewBanService(repo *repository.PlayerRepository) *BanService {
	return &BanService{repo: repo}
}

type BanRequest struct {
	Reason    string
	Duration  int // Duration in hours (0 = permanent)
	ExpiresAt *time.Time
}

// BanPlayer bans a player with optional expiration
func (s *BanService) BanPlayer(ctx context.Context, playerID string, req BanRequest) error {
	player, err := s.repo.GetPlayerByID(ctx, playerID)
	if err != nil {
		return err
	}

	if player == nil {
		return fmt.Errorf("player not found")
	}

	// Check if player is owner (cannot ban owners)
	for _, role := range player.Role {
		if role == models.ROLE_OWNER {
			return fmt.Errorf("cannot ban owner")
		}
	}

	now := time.Now()
	player.IsBanned = true
	player.BannedAt = &now
	player.BanReason = req.Reason

	// Set expiration
	if req.ExpiresAt != nil {
		player.BanExpiresAt = req.ExpiresAt
	} else if req.Duration > 0 {
		expiry := now.Add(time.Duration(req.Duration) * time.Hour)
		player.BanExpiresAt = &expiry
	} else {
		// Permanent ban - no expiry
		player.BanExpiresAt = nil
	}

	return s.repo.UpdatePlayer(ctx, player)
}

// UnbanPlayer removes ban from a player
func (s *BanService) UnbanPlayer(ctx context.Context, playerID string) error {
	player, err := s.repo.GetPlayerByID(ctx, playerID)
	if err != nil {
		return err
	}

	if player == nil {
		return fmt.Errorf("player not found")
	}

	player.IsBanned = false
	player.BanReason = ""
	player.BanExpiresAt = nil
	player.BannedAt = nil

	return s.repo.UpdatePlayer(ctx, player)
}

// CheckExpiredBans finds and unbans all players with expired bans
func (s *BanService) CheckExpiredBans(ctx context.Context) (int64, error) {
	return s.repo.UnbanExpiredPlayers(ctx)
}

// GetBanInfo returns ban information for a player
type BanInfo struct {
	IsBanned    bool       `json:"is_banned"`
	Reason      string     `json:"reason"`
	BannedAt    *time.Time `json:"banned_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsExpired   bool       `json:"is_expired"`
	TimeLeft    string     `json:"time_left"` // Human readable remaining time
	IsPermanent bool       `json:"is_permanent"`
}

func (s *BanService) GetBanInfo(ctx context.Context, playerID string) (*BanInfo, error) {
	player, err := s.repo.GetPlayerByID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	if player == nil {
		return nil, fmt.Errorf("player not found")
	}

	info := &BanInfo{
		IsBanned:    player.IsBanned,
		Reason:      player.BanReason,
		BannedAt:    player.BannedAt,
		ExpiresAt:   player.BanExpiresAt,
		IsPermanent: player.BanExpiresAt == nil && player.IsBanned,
	}

	if player.IsBanned && player.BanExpiresAt != nil {
		now := time.Now()
		if now.After(*player.BanExpiresAt) {
			info.IsExpired = true
			info.TimeLeft = "Expired"
		} else {
			timeLeft := player.BanExpiresAt.Sub(now)
			info.TimeLeft = formatDuration(timeLeft)
		}
	}

	return info, nil
}

// GetBannedPlayers returns all currently banned players
func (s *BanService) GetBannedPlayers(ctx context.Context) ([]models.Player, error) {
	return s.repo.GetBannedPlayers(ctx)
}

// GetExpiringSoonBans returns players whose bans are expiring soon (within 24 hours)
func (s *BanService) GetExpiringSoonBans(ctx context.Context) ([]models.Player, error) {
	return s.repo.GetExpiringSoonBans(ctx, 24*time.Hour)
}

// Helper function to format duration in human-readable format
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
