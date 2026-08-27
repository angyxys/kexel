package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type InvitationHandler struct {
	serv *service.InvitationService
}

func NewInvitationHandler(serv *service.InvitationService) *InvitationHandler {
	return &InvitationHandler{
		serv: serv,
	}
}

type CreateInvitationRequest struct {
	Role      string `json:"role" binding:"required"`
	MaxUses   int    `json:"max_uses" binding:"min=-1"`
	ExpiresAt string `json:"expires_at"` // ISO 8601 format
}

type InvitationResponse struct {
	ID        uint      `json:"id"`
	Code      string    `json:"code"`
	Role      string    `json:"role"`
	MaxUses   int       `json:"max_uses"`
	Uses      int       `json:"uses"`
	ExpiresAt *time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateInvitation generates a new invitation code
func (h *InvitationHandler) CreateInvitation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")

	var req CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	invitationReq := service.InvitationRequest{
		Role:    req.Role,
		MaxUses: req.MaxUses,
	}

	// Parse expiration date if provided
	if req.ExpiresAt != "" {
		if expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt); err == nil {
			invitationReq.ExpiresAt = &expiresAt
		}
	}

	invitation, err := h.serv.CreateInvitation(ctx, userID.(uint), invitationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, InvitationResponse{
		ID:        invitation.ID,
		Code:      invitation.Code,
		Role:      invitation.Role,
		MaxUses:   invitation.MaxUses,
		Uses:      invitation.Uses,
		ExpiresAt: invitation.ExpiresAt,
		IsActive:  invitation.IsActive,
		CreatedAt: invitation.CreatedAt,
	})
}

// GetMyInvitations retrieves all invitations created by the current user
func (h *InvitationHandler) GetMyInvitations(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")

	invitations, err := h.serv.GetInvitations(ctx, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve invitations",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	responses := make([]InvitationResponse, len(invitations))
	for i, inv := range invitations {
		responses[i] = InvitationResponse{
			ID:        inv.ID,
			Code:      inv.Code,
			Role:      inv.Role,
			MaxUses:   inv.MaxUses,
			Uses:      inv.Uses,
			ExpiresAt: inv.ExpiresAt,
			IsActive:  inv.IsActive,
			CreatedAt: inv.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// ValidateInvitation checks if an invitation code is valid
func (h *InvitationHandler) ValidateInvitation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "code parameter required",
			"status":  http.StatusBadRequest,
		})
		return
	}

	invitation, err := h.serv.ValidateInvitation(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"role":    invitation.Role,
		"uses":    invitation.Uses,
		"max_uses": invitation.MaxUses,
		"status":  http.StatusOK,
	})
}

// RevokeInvitation deactivates an invitation code
func (h *InvitationHandler) RevokeInvitation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid invitation id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.serv.RevokeInvitation(ctx, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Invitation revoked successfully",
		"status":  http.StatusOK,
	})
}

// GetInvitationStats retrieves statistics about active invitations
func (h *InvitationHandler) GetInvitationStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	invitations, err := h.serv.GetActiveInvitations(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve stats",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	totalUses := 0
	totalCapacity := 0

	for _, inv := range invitations {
		totalUses += inv.Uses
		if inv.MaxUses != -1 {
			totalCapacity += inv.MaxUses
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_active":      len(invitations),
		"total_uses":        totalUses,
		"total_capacity":    totalCapacity,
		"status":            http.StatusOK,
	})
}
