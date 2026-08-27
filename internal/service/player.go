package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type PlayerService struct {
	repo *repository.PlayerRepository
}

func NewPlayerService(repo *repository.PlayerRepository) *PlayerService {
	return &PlayerService{
		repo: repo,
	}
}

func (s *PlayerService) GetPlayerByID(ctx context.Context, playerID string) (*models.Player, error) {
	return s.repo.GetPlayerByID(ctx, playerID)
}

func (s *PlayerService) AddPlayer(ctx context.Context, player *models.Player) error {
	exists, err := s.repo.GetPlayerByID(ctx, player.VRChatID)
	if err != nil {
		return err
	}
	if exists != nil {
		return fmt.Errorf("player already exists")
	}
	return s.repo.AddPlayer(ctx, player)
}

func (s *PlayerService) UpdatePlayer(ctx context.Context, player *models.Player) error {
	return s.repo.UpdatePlayer(ctx, player)
}

func (s *PlayerService) DeletePlayer(ctx context.Context, playerID string) error {
	return s.repo.DeletePlayer(ctx, playerID)
}

func (s *PlayerService) ListPlayers(ctx context.Context) ([]models.Player, error) {
	return s.repo.ListPlayers(ctx)
}

func (s *PlayerService) SearchPlayers(ctx context.Context, query string, page int, pageSize int) ([]models.Player, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	return s.repo.SearchPlayers(ctx, query, pageSize, offset)
}

func (s *PlayerService) FilterPlayers(ctx context.Context, filters repository.PlayerFilters, page int, pageSize int) ([]models.Player, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	return s.repo.FilterPlayers(ctx, filters, pageSize, offset)
}

func (s *PlayerService) BanUnbanPlayer(ctx context.Context, playerID string, ban bool) error {
	exists, err := s.repo.GetPlayerByID(ctx, playerID)
	if err != nil {
		return err
	}
	if exists == nil {
		return fmt.Errorf("player not found")
	}
	if slices.Contains(exists.Role, models.ROLE_OWNER) {
		return fmt.Errorf("cannot ban owner")
	}
	return s.repo.BanUnbanPlayer(ctx, playerID, ban)
}

func (s *PlayerService) HasRoles(ctx context.Context, playerID string, roles []models.Role) bool {
	player, err := s.repo.GetPlayerByID(ctx, playerID)
	if err != nil || player == nil {
		return false
	}
	playerRolesMap := make(map[models.Role]struct{}, len(player.Role))
	for _, role := range player.Role {
		playerRolesMap[role] = struct{}{}
	}
	for _, reqRole := range roles {
		if _, hasRole := playerRolesMap[reqRole]; !hasRole {
			return false
		}
	}
	return true
}

func (s *PlayerService) AllVip(ctx context.Context) []models.Player {
	return s.repo.AllVip(ctx)
}

func (s *PlayerService) AllBanned(ctx context.Context) []models.Player {
	return s.repo.AllBanned(ctx)
}

func (s *PlayerService) AllModerator(ctx context.Context) []models.Player {
	return s.repo.AllModerator(ctx)
}

func (s *PlayerService) AllOwner(ctx context.Context) []models.Player {
	return s.repo.AllOwner(ctx)
}
