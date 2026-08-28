package handler

import (
	"net/http"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type IntegrationHandler struct {
	discordServ   *service.DiscordService
	patreonServ   *service.PatreonService
	rateLimitServ *service.RateLimitService
}

func NewIntegrationHandler(discordServ *service.DiscordService, patreonServ *service.PatreonService, rateLimitServ *service.RateLimitService) *IntegrationHandler {
	return &IntegrationHandler{
		discordServ:   discordServ,
		patreonServ:   patreonServ,
		rateLimitServ: rateLimitServ,
	}
}

type SetupDiscordRequest struct {
	BotToken string `json:"bot_token" binding:"required"`
	GuildID  string `json:"guild_id" binding:"required"`
}

type ConfigureChannelsRequest struct {
	ModLogChannelID       string `json:"mod_log_channel_id"`
	AnnouncementChannelID string `json:"announcement_channel_id"`
}

// SetupDiscord sets up Discord bot integration
func (h *IntegrationHandler) SetupDiscord(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req SetupDiscordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	discord, err := h.discordServ.SetupDiscord(c, userID.(uint), req.BotToken, req.GuildID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, discord)
}

// GetDiscordIntegration gets Discord integration info
func (h *IntegrationHandler) GetDiscordIntegration(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	discord, err := h.discordServ.GetDiscordIntegration(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"status":  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, discord)
}

// ConfigureDiscordChannels configures Discord channels
func (h *IntegrationHandler) ConfigureDiscordChannels(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req ConfigureChannelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.discordServ.ConfigureChannels(c, userID.(uint), req.ModLogChannelID, req.AnnouncementChannelID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "channels configured",
		"status":  http.StatusOK,
	})
}

// TestDiscordConnection tests Discord bot connection
func (h *IntegrationHandler) TestDiscordConnection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	if err := h.discordServ.TestConnection(c, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "connection successful",
		"status":  http.StatusOK,
	})
}

// DisconnectDiscord disconnects Discord integration
func (h *IntegrationHandler) DisconnectDiscord(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	if err := h.discordServ.DisconnectDiscord(c, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "discord disconnected",
		"status":  http.StatusOK,
	})
}

// GetPatreonOAuthURL gets the Patreon OAuth URL
func (h *IntegrationHandler) GetPatreonOAuthURL(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "state parameter required",
			"status":  http.StatusBadRequest,
		})
		return
	}

	url := h.patreonServ.GetOAuthURL(state)
	c.JSON(http.StatusOK, gin.H{
		"url":    url,
		"status": http.StatusOK,
	})
}

type PatreonOAuthCallbackRequest struct {
	Code string `json:"code" binding:"required"`
}

// HandlePatreonOAuthCallback handles Patreon OAuth callback
func (h *IntegrationHandler) HandlePatreonOAuthCallback(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req PatreonOAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	patreon, err := h.patreonServ.HandleOAuthCallback(c, userID.(uint), req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, patreon)
}

// GetPatreonIntegration gets Patreon integration info
func (h *IntegrationHandler) GetPatreonIntegration(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	patreon, err := h.patreonServ.GetPatreonInfo(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"status":  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, patreon)
}

type ConfigureTierMappingRequest struct {
	TierMapping map[string]string `json:"tier_mapping" binding:"required"`
}

// ConfigurePatreonTierMapping configures tier-to-role mapping
func (h *IntegrationHandler) ConfigurePatreonTierMapping(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req ConfigureTierMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.patreonServ.ConfigureTierMapping(c, userID.(uint), req.TierMapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "tier mapping configured",
		"status":  http.StatusOK,
	})
}

// SyncPatreonMembers syncs Patreon members
func (h *IntegrationHandler) SyncPatreonMembers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	if err := h.patreonServ.SyncPatreons(c, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"status":  http.StatusOK,
	})
}

// DisconnectPatreon disconnects Patreon integration
func (h *IntegrationHandler) DisconnectPatreon(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	if err := h.patreonServ.DisconnectPatreon(c, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "patreon disconnected",
		"status":  http.StatusOK,
	})
}

// GetRateLimitBlocks gets recent rate limit blocks
func (h *IntegrationHandler) GetRateLimitBlocks(c *gin.Context) {
	blocks, err := h.rateLimitServ.GetRecentBlocks(c, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch blocks",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   blocks,
		"status": http.StatusOK,
	})
}
