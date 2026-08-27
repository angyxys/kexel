package service

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type StatsService struct {
	playerRepo *repository.PlayerRepository
	auditRepo  *repository.AuditLogRepository
	userRepo   *repository.UserRepository
}

func NewStatsService(playerRepo *repository.PlayerRepository, auditRepo *repository.AuditLogRepository, userRepo *repository.UserRepository) *StatsService {
	return &StatsService{
		playerRepo: playerRepo,
		auditRepo:  auditRepo,
		userRepo:   userRepo,
	}
}

// KPIStats contains key performance indicators
type KPIStats struct {
	TotalPlayers  int64 `json:"total_players"`
	TotalVIPs     int64 `json:"total_vips"`
	TotalMods     int64 `json:"total_mods"`
	TotalOwners   int64 `json:"total_owners"`
	TotalBanned   int64 `json:"total_banned"`
	TotalActive   int64 `json:"total_active"`
	TotalUsers    int64 `json:"total_users"`
	TodayLogins   int64 `json:"today_logins"`
	BansThisWeek  int64 `json:"bans_this_week"`
	UnbansThisWeek int64 `json:"unbans_this_week"`
}

// GetKPIStats retrieves key performance indicators
func (s *StatsService) GetKPIStats(ctx context.Context) (*KPIStats, error) {
	stats := &KPIStats{}

	// Get total players
	players, err := s.playerRepo.ListPlayers(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalPlayers = int64(len(players))

	// Count by role and ban status
	vipCount := 0
	modCount := 0
	ownerCount := 0
	bannedCount := 0

	for _, player := range players {
		if player.IsBanned {
			bannedCount++
		} else {
			stats.TotalActive++
		}

		for _, role := range player.Role {
			switch role {
			case models.ROLE_VIP:
				vipCount++
			case models.ROLE_MOD:
				modCount++
			case models.ROLE_OWNER:
				ownerCount++
			}
		}
	}

	stats.TotalVIPs = int64(vipCount)
	stats.TotalMods = int64(modCount)
	stats.TotalOwners = int64(ownerCount)
	stats.TotalBanned = int64(bannedCount)
	stats.TotalActive = stats.TotalPlayers - int64(bannedCount)

	// Get total registered users
	users, err := s.userRepo.ListAllUsers(ctx)
	if err == nil {
		stats.TotalUsers = int64(len(users))
	}

	// Get today's logins (from audit logs)
	filters := map[string]interface{}{
		"action":     "POST",
		"start_date": time.Now().Add(-24 * time.Hour),
	}
	logs, _, err := s.auditRepo.List(ctx, filters, 1000, 0)
	if err == nil {
		stats.TodayLogins = int64(len(logs))
	}

	// Count bans and unbans this week
	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
	for _, log := range logs {
		if log.CreatedAt.After(weekAgo) {
			if log.Action == "PUT" && log.ResourceType == "player" {
				// Could be a ban or unban, need to check old_value/new_value
				// For now, count as potential ban action
			}
		}
	}

	return stats, nil
}

// ActivityStats contains recent activity data
type ActivityLog struct {
	ID        uint      `json:"id"`
	User      string    `json:"user"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
}

// GetRecentActivity retrieves recent activities
func (s *StatsService) GetRecentActivity(ctx context.Context, limit int) ([]ActivityLog, error) {
	if limit > 50 {
		limit = 50
	}

	filters := make(map[string]interface{})
	logs, _, err := s.auditRepo.List(ctx, filters, limit, 0)
	if err != nil {
		return nil, err
	}

	activities := make([]ActivityLog, len(logs))
	for i, log := range logs {
		activities[i] = ActivityLog{
			ID:        log.ID,
			User:      log.User.Username,
			Action:    log.Action,
			Resource:  log.ResourceType + ": " + log.ResourceID,
			Timestamp: log.CreatedAt,
			Details:   log.Description,
		}
	}

	return activities, nil
}

// TrendData contains trend information
type TrendData struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// GetPlayerTrends retrieves player count trends over time
func (s *StatsService) GetPlayerTrends(ctx context.Context, days int) ([]TrendData, error) {
	trends := make([]TrendData, 0)

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		// In a real scenario, this would query historical data
		// For now, return approximate data based on current state
		trends = append(trends, TrendData{
			Date:  date,
			Value: int64(i * 10), // Mock data
		})
	}

	return trends, nil
}

// GetRoleDistribution returns distribution of players by role
type RoleDistribution struct {
	Role  string `json:"role"`
	Count int64  `json:"count"`
}

func (s *StatsService) GetRoleDistribution(ctx context.Context) ([]RoleDistribution, error) {
	players, err := s.playerRepo.ListPlayers(ctx)
	if err != nil {
		return nil, err
	}

	distribution := make(map[string]int64)
	distribution["user"] = 0
	distribution["vip"] = 0
	distribution["mod"] = 0
	distribution["owner"] = 0

	for _, player := range players {
		for _, role := range player.Role {
			distribution[string(role)]++
		}
	}

	result := make([]RoleDistribution, 0)
	for role, count := range distribution {
		result = append(result, RoleDistribution{
			Role:  role,
			Count: count,
		})
	}

	return result, nil
}

// GetBanStats returns ban-related statistics
type BanStats struct {
	TotalBanned     int64 `json:"total_banned"`
	PermanentBans   int64 `json:"permanent_bans"`
	TemporaryBans   int64 `json:"temporary_bans"`
	ExpiringToday   int64 `json:"expiring_today"`
	ExpiringWeek    int64 `json:"expiring_week"`
	MostCommonReason string `json:"most_common_reason"`
}

func (s *StatsService) GetBanStats(ctx context.Context) (*BanStats, error) {
	players, err := s.playerRepo.GetBannedPlayers(ctx)
	if err != nil {
		return nil, err
	}

	stats := &BanStats{
		TotalBanned: int64(len(players)),
	}

	now := time.Now()
	todayEnd := now.AddDate(0, 0, 1)
	weekEnd := now.AddDate(0, 0, 7)

	reasonCounts := make(map[string]int)

	for _, player := range players {
		if player.BanExpiresAt == nil {
			stats.PermanentBans++
		} else {
			stats.TemporaryBans++

			if player.BanExpiresAt.Before(todayEnd) && player.BanExpiresAt.After(now) {
				stats.ExpiringToday++
			}
			if player.BanExpiresAt.Before(weekEnd) && player.BanExpiresAt.After(now) {
				stats.ExpiringWeek++
			}
		}

		if player.BanReason != "" {
			reasonCounts[player.BanReason]++
		}
	}

	// Find most common reason
	maxCount := 0
	for reason, count := range reasonCounts {
		if count > maxCount {
			maxCount = count
			stats.MostCommonReason = reason
		}
	}

	return stats, nil
}
