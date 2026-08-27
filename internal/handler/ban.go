package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type BanHandler struct {
	serv *service.BanService
}

func NewBanHandler(serv *service.BanService) *BanHandler {
	return &BanHandler{
		serv: serv,
	}
}

type BanPlayerRequest struct {
	Reason    string `json:"reason" binding:"required"`
	Duration  int    `json:"duration"` // Duration in hours (0 = permanent)
	ExpiresAt string `json:"expires_at"` // ISO 8601 date format
}

type UnbanPlayerRequest struct {
	Reason string `json:"reason"`
}

type BanInfoResponse struct {
	VRChatID    string `json:"vrchat_id"`
	IsBanned    bool   `json:"is_banned"`
	Reason      string `json:"reason"`
	BannedAt    *time.Time `json:"banned_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsExpired   bool   `json:"is_expired"`
	TimeLeft    string `json:"time_left"`
	IsPermanent bool   `json:"is_permanent"`
}

// BanPlayer bans a player with optional expiration
func (h *BanHandler) BanPlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playerID := c.Param("id")

	var req BanPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	banReq := service.BanRequest{
		Reason:   req.Reason,
		Duration: req.Duration,
	}

	// Parse expiration date if provided
	if req.ExpiresAt != "" {
		if expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt); err == nil {
			banReq.ExpiresAt = &expiresAt
		}
	}

	err := h.serv.BanPlayer(ctx, playerID, banReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Player banned successfully",
		"status":  http.StatusOK,
	})
}

// UnbanPlayer removes ban from a player
func (h *BanHandler) UnbanPlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playerID := c.Param("id")

	err := h.serv.UnbanPlayer(ctx, playerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Player unbanned successfully",
		"status":  http.StatusOK,
	})
}

// GetBanInfo returns ban information for a player
func (h *BanHandler) GetBanInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playerID := c.Param("id")

	info, err := h.serv.GetBanInfo(ctx, playerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"status":  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetBannedPlayers returns all currently banned players
func (h *BanHandler) GetBannedPlayers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	players, err := h.serv.GetBannedPlayers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve banned players",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	responses := make([]BanInfoResponse, 0)
	for _, player := range players {
		info, _ := h.serv.GetBanInfo(ctx, player.VRChatID)
		if info != nil {
			responses = append(responses, BanInfoResponse{
				VRChatID:    player.VRChatID,
				IsBanned:    info.IsBanned,
				Reason:      info.Reason,
				BannedAt:    info.BannedAt,
				ExpiresAt:   info.ExpiresAt,
				IsExpired:   info.IsExpired,
				TimeLeft:    info.TimeLeft,
				IsPermanent: info.IsPermanent,
			})
		}
	}

	c.JSON(http.StatusOK, responses)
}

// GetExpiringSoonBans returns bans expiring within 24 hours
func (h *BanHandler) GetExpiringSoonBans(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	players, err := h.serv.GetExpiringSoonBans(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve expiring bans",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	responses := make([]BanInfoResponse, 0)
	for _, player := range players {
		info, _ := h.serv.GetBanInfo(ctx, player.VRChatID)
		if info != nil {
			responses = append(responses, BanInfoResponse{
				VRChatID:    player.VRChatID,
				IsBanned:    info.IsBanned,
				Reason:      info.Reason,
				BannedAt:    info.BannedAt,
				ExpiresAt:   info.ExpiresAt,
				IsExpired:   info.IsExpired,
				TimeLeft:    info.TimeLeft,
				IsPermanent: info.IsPermanent,
			})
		}
	}

	c.JSON(http.StatusOK, responses)
}

// CleanupExpiredBans manually triggers cleanup of expired bans
func (h *BanHandler) CleanupExpiredBans(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := h.serv.CheckExpiredBans(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cleanup expired bans",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cleanup completed",
		"unbanned_count": count,
		"status":  http.StatusOK,
	})
}
