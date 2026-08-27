package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	serv *service.PlayerService
}

func NewPlayerHandler(serv *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		serv: serv,
	}
}

type PlayerRequest struct {
	VRChatID string         `json:"vrchat_id" binding:"required"`
	Roles    *[]models.Role `json:"roles"`
	IsBanned *bool          `json:"is_banned"`
}

type PlayerResponse struct {
	VRChatID string         `json:"vrchat_id"`
	Role     []models.Role  `json:"roles"`
	IsBanned bool           `json:"is_banned"`
}

func (h *PlayerHandler) Vip(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	c.JSON(http.StatusOK, h.serv.AllVip(ctx))
}

func (h *PlayerHandler) Banned(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	c.JSON(http.StatusOK, h.serv.AllBanned(ctx))
}

func (h *PlayerHandler) Moderator(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	c.JSON(http.StatusOK, h.serv.AllModerator(ctx))
}

func (h *PlayerHandler) Owner(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	c.JSON(http.StatusOK, h.serv.AllOwner(ctx))
}

// AddPlayer creates a new player
func (h *PlayerHandler) AddPlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userRole, _ := c.Get("role")
	// Only owner can create players
	if userRole != string(models.ROLE_OWNER) && userRole != string(models.ROLE_MOD) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "insufficient permissions",
			"status":  http.StatusForbidden,
		})
		return
	}

	var playerReq PlayerRequest
	if err := c.ShouldBindJSON(&playerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	player := &models.Player{
		VRChatID: playerReq.VRChatID,
	}

	if playerReq.Roles != nil {
		player.Role = *playerReq.Roles
	}
	if playerReq.IsBanned != nil {
		player.IsBanned = *playerReq.IsBanned
	}

	if err := h.serv.AddPlayer(ctx, player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("player %s created successfully", playerReq.VRChatID),
		"status":  http.StatusCreated,
		"player": PlayerResponse{
			VRChatID: player.VRChatID,
			Role:     player.Role,
			IsBanned: player.IsBanned,
		},
	})
}

// GetPlayer retrieves a player by ID
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	vrchatID := c.Param("id")

	player, err := h.serv.GetPlayerByID(ctx, vrchatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	if player == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "player not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, PlayerResponse{
		VRChatID: player.VRChatID,
		Role:     player.Role,
		IsBanned: player.IsBanned,
	})
}

// UpdatePlayer updates a player
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userRole, _ := c.Get("role")
	// Only owner can update players
	if userRole != string(models.ROLE_OWNER) && userRole != string(models.ROLE_MOD) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "insufficient permissions",
			"status":  http.StatusForbidden,
		})
		return
	}

	vrchatID := c.Param("id")

	var playerReq PlayerRequest
	if err := c.ShouldBindJSON(&playerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	player := &models.Player{
		VRChatID: vrchatID,
	}

	if playerReq.Roles != nil {
		player.Role = *playerReq.Roles
	}
	if playerReq.IsBanned != nil {
		player.IsBanned = *playerReq.IsBanned
	}

	if err := h.serv.UpdatePlayer(ctx, player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "player updated successfully",
		"status":  http.StatusOK,
		"player": PlayerResponse{
			VRChatID: player.VRChatID,
			Role:     player.Role,
			IsBanned: player.IsBanned,
		},
	})
}

// DeletePlayer deletes a player
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userRole, _ := c.Get("role")
	// Only owner can delete players
	if userRole != string(models.ROLE_OWNER) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "insufficient permissions",
			"status":  http.StatusForbidden,
		})
		return
	}

	vrchatID := c.Param("id")

	if err := h.serv.DeletePlayer(ctx, vrchatID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "player deleted successfully",
		"status":  http.StatusOK,
	})
}

// ListPlayers retrieves all players
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	players, err := h.serv.ListPlayers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	playerResponses := make([]PlayerResponse, len(players))
	for i, p := range players {
		playerResponses[i] = PlayerResponse{
			VRChatID: p.VRChatID,
			Role:     p.Role,
			IsBanned: p.IsBanned,
		}
	}

	c.JSON(http.StatusOK, playerResponses)
}

// SearchPlayers performs full-text search on players
func (h *PlayerHandler) SearchPlayers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "search query required",
			"status":  http.StatusBadRequest,
		})
		return
	}

	players, total, err := h.serv.SearchPlayers(ctx, query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	playerResponses := make([]PlayerResponse, len(players))
	for i, p := range players {
		playerResponses[i] = PlayerResponse{
			VRChatID: p.VRChatID,
			Role:     p.Role,
			IsBanned: p.IsBanned,
		}
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"data":        playerResponses,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// FilterPlayers applies multiple filters to player list
func (h *PlayerHandler) FilterPlayers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := repository.PlayerFilters{
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "vrchat_id"),
		SortOrder: c.DefaultQuery("sort_order", "ASC"),
	}

	// Parse roles filter
	if rolesStr := c.Query("roles"); rolesStr != "" {
		filters.Roles = strings.Split(rolesStr, ",")
	}

	// Parse ban status filter
	if bannedStr := c.Query("banned"); bannedStr != "" {
		banned := bannedStr == "true"
		filters.IsBanned = &banned
	}

	players, total, err := h.serv.FilterPlayers(ctx, filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	playerResponses := make([]PlayerResponse, len(players))
	for i, p := range players {
		playerResponses[i] = PlayerResponse{
			VRChatID: p.VRChatID,
			Role:     p.Role,
			IsBanned: p.IsBanned,
		}
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"data":        playerResponses,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}
