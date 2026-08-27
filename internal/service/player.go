package service

import (
	"context"
	"log"
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

func (s *PlayerService) GetPlayerByID(ctx context.Context, playerID string) *models.Player {
	return s.repo.GetPlayerByID(ctx, playerID)
}

func (s *PlayerService) AddPlayer(ctx context.Context, player *models.Player) {
	exists := s.repo.GetPlayerByID(ctx, player.VRChatID)
	if exists != nil {
		log.Println("user already exits")
		return
	}
	s.repo.AddPlayer(ctx, player)
}

func (s *PlayerService) BanUnbanPlayer(ctx context.Context, playerID string, ban bool) {
	exists := s.repo.GetPlayerByID(ctx, playerID)
	if exists == nil {
		log.Println("user not exists")
		return
	}
	if slices.Contains(exists.Role, models.ROLE_OWNER) {
		log.Println("cannot ban owner")
		return
	}
	s.repo.BanUnbanPlayer(ctx, playerID, ban)
}

func (s *PlayerService) HasRoles(ctx context.Context, playerID string, roles []models.Role) bool {
	player := s.repo.GetPlayerByID(ctx, playerID)
	if player == nil {
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
