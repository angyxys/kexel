package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
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

func (h *PlayerHandler) AddPlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	playerID := c.GetString("playerID")
	hasRoles := h.serv.HasRoles(ctx, playerID, []models.Role{
		models.ROLE_OWNER,
	})
	if !hasRoles {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you cannot permission to execute this action",
			"status":  http.StatusForbidden,
		})
		return
	}
	var player struct {
		VRChatID string         `json:"vrchat_id" binding:"required"`
		Roles    *[]models.Role `json:"roles"`
		IsBanned *bool          `json:"is_banned"`
	}
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}
	insertPlayer := &models.Player{
		VRChatID: player.VRChatID,
	}
	if player.Roles != nil {
		insertPlayer.Role = *player.Roles
	}
	if player.IsBanned != nil {
		insertPlayer.IsBanned = *player.IsBanned
	}
	h.serv.AddPlayer(ctx, insertPlayer)
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("user %s created successfully", player.VRChatID),
		"status":  http.StatusCreated,
	})
}
