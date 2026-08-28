package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type APIKeyHandler struct {
	apiKeyServ *service.APIKeyService
}

func NewAPIKeyHandler(apiKeyServ *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyServ: apiKeyServ,
	}
}

type CreateAPIKeyRequest struct {
	Name      string    `json:"name" binding:"required,min=1,max=100"`
	Scopes    []string  `json:"scopes" binding:"required,min=1"`
	ExpiresIn *int      `json:"expires_in"` // Days until expiry (optional)
	RateLimit int       `json:"rate_limit" binding:"min=1"`
}

type UpdateRateLimitRequest struct {
	RateLimit int `json:"rate_limit" binding:"required,min=1"`
}

type UpdateScopesRequest struct {
	Scopes []string `json:"scopes" binding:"required,min=1"`
}

// CreateAPIKey creates a new API key
func (h *APIKeyHandler) CreateAPIKey(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		exp := time.Now().AddDate(0, 0, *req.ExpiresIn)
		expiresAt = &exp
	}

	if req.RateLimit == 0 {
		req.RateLimit = 1000
	}

	apiKey, err := h.apiKeyServ.CreateAPIKey(c, userID.(uint), req.Name, req.Scopes, expiresAt, req.RateLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, apiKey)
}

// GetAPIKeys returns all API keys for the authenticated user
func (h *APIKeyHandler) GetAPIKeys(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	keys, err := h.apiKeyServ.GetUserAPIKeys(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch API keys",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   keys,
		"status": http.StatusOK,
	})
}

// RevokeAPIKey revokes an API key
func (h *APIKeyHandler) RevokeAPIKey(c *gin.Context) {
	keyIDStr := c.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid key id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.apiKeyServ.RevokeAPIKey(c, uint(keyID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to revoke API key",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key revoked",
		"status":  http.StatusOK,
	})
}

// DeleteAPIKey deletes an API key
func (h *APIKeyHandler) DeleteAPIKey(c *gin.Context) {
	keyIDStr := c.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid key id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.apiKeyServ.DeleteAPIKey(c, uint(keyID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete API key",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key deleted",
		"status":  http.StatusOK,
	})
}

// UpdateRateLimit updates the rate limit for an API key
func (h *APIKeyHandler) UpdateRateLimit(c *gin.Context) {
	keyIDStr := c.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid key id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	var req UpdateRateLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.apiKeyServ.UpdateRateLimit(c, uint(keyID), req.RateLimit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "rate limit updated",
		"status":  http.StatusOK,
	})
}

// UpdateScopes updates the scopes for an API key
func (h *APIKeyHandler) UpdateScopes(c *gin.Context) {
	keyIDStr := c.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid key id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	var req UpdateScopesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.apiKeyServ.UpdateScopes(c, uint(keyID), req.Scopes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "scopes updated",
		"status":  http.StatusOK,
	})
}

// GetAvailableScopes returns all available API scopes
func (h *APIKeyHandler) GetAvailableScopes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"scopes": service.AvailableScopes(),
		"status": http.StatusOK,
	})
}
