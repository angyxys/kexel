package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type DiscordService struct {
	discordRepo *repository.DiscordIntegrationRepository
}

func NewDiscordService(discordRepo *repository.DiscordIntegrationRepository) *DiscordService {
	return &DiscordService{
		discordRepo: discordRepo,
	}
}

// SetupDiscord sets up Discord bot integration
func (s *DiscordService) SetupDiscord(ctx context.Context, userID uint, botToken string, guildID string) (*DiscordInfo, error) {
	if botToken == "" || guildID == "" {
		return nil, errors.New("bot token and guild ID are required")
	}

	// Check if integration already exists
	existing, err := s.discordRepo.GetByUserID(ctx, userID)
	if err == nil && existing != nil {
		// Update existing
		existing.BotToken = botToken
		existing.GuildID = guildID
		if err := s.discordRepo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return s.getDiscordInfo(existing), nil
	}

	// Create new integration
	integration := &models.DiscordIntegration{
		UserID:   userID,
		BotToken: botToken,
		GuildID:  guildID,
	}

	if err := s.discordRepo.Create(ctx, integration); err != nil {
		return nil, fmt.Errorf("error setting up Discord: %w", err)
	}

	return s.getDiscordInfo(integration), nil
}

type DiscordInfo struct {
	ID                    uint   `json:"id"`
	GuildID               string `json:"guild_id"`
	ModLogChannelID       string `json:"mod_log_channel_id"`
	AnnouncementChannelID string `json:"announcement_channel_id"`
	IsConnected           bool   `json:"is_connected"`
}

func (s *DiscordService) getDiscordInfo(integration *models.DiscordIntegration) *DiscordInfo {
	return &DiscordInfo{
		ID:                    integration.ID,
		GuildID:               integration.GuildID,
		ModLogChannelID:       integration.ModLogChannelID,
		AnnouncementChannelID: integration.AnnouncementChannelID,
		IsConnected:           integration.IsConnected,
	}
}

// GetDiscordIntegration returns Discord integration info
func (s *DiscordService) GetDiscordIntegration(ctx context.Context, userID uint) (*DiscordInfo, error) {
	integration, err := s.discordRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("Discord not configured")
	}

	return s.getDiscordInfo(integration), nil
}

// ConfigureChannels configures Discord channels
func (s *DiscordService) ConfigureChannels(ctx context.Context, userID uint, modLogChannel string, announcementChannel string) error {
	integration, err := s.discordRepo.GetByUserID(ctx, userID)
	if err != nil {
		return errors.New("Discord not configured")
	}

	integration.ModLogChannelID = modLogChannel
	integration.AnnouncementChannelID = announcementChannel

	return s.discordRepo.Update(ctx, integration)
}

// ConfigureRoleMapping configures role mapping between Kexel and Discord
func (s *DiscordService) ConfigureRoleMapping(ctx context.Context, userID uint, mapping map[string]string) error {
	integration, err := s.discordRepo.GetByUserID(ctx, userID)
	if err != nil {
		return errors.New("Discord not configured")
	}

	roleJSON, _ := json.Marshal(mapping)
	integration.RoleMapping = string(roleJSON)

	return s.discordRepo.Update(ctx, integration)
}

// DisconnectDiscord disconnects Discord integration
func (s *DiscordService) DisconnectDiscord(ctx context.Context, userID uint) error {
	return s.discordRepo.Delete(ctx, userID)
}

// TestConnection tests Discord bot connection
func (s *DiscordService) TestConnection(ctx context.Context, userID uint) error {
	_, err := s.discordRepo.GetByUserID(ctx, userID)
	if err != nil {
		return errors.New("Discord not configured")
	}

	// In production, this would actually test the Discord bot connection
	// For now, just mark as connected
	return s.discordRepo.UpdateConnectionStatus(ctx, userID, true)
}
