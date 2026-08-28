package handler

import (
	"net/http"
	"strconv"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	webhookServ *service.WebhookService
}

func NewWebhookHandler(webhookServ *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookServ: webhookServ,
	}
}

type CreateWebhookRequest struct {
	Name   string   `json:"name" binding:"required,min=1,max=100"`
	URL    string   `json:"url" binding:"required,url"`
	Events []string `json:"events" binding:"required,min=1"`
}

// CreateWebhook creates a new webhook
func (h *WebhookHandler) CreateWebhook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	webhook, err := h.webhookServ.CreateWebhook(c, userID.(uint), req.Name, req.URL, req.Events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, webhook)
}

// GetWebhooks returns all webhooks for the authenticated user
func (h *WebhookHandler) GetWebhooks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	webhooks, err := h.webhookServ.GetUserWebhooks(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch webhooks",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   webhooks,
		"status": http.StatusOK,
	})
}

// DeleteWebhook deletes a webhook
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	webhookIDStr := c.Param("id")
	webhookID, err := strconv.ParseUint(webhookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid webhook id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.webhookServ.DeleteWebhook(c, uint(webhookID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete webhook",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "webhook deleted",
		"status":  http.StatusOK,
	})
}

// DisableWebhook disables a webhook
func (h *WebhookHandler) DisableWebhook(c *gin.Context) {
	webhookIDStr := c.Param("id")
	webhookID, err := strconv.ParseUint(webhookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid webhook id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.webhookServ.DisableWebhook(c, uint(webhookID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to disable webhook",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "webhook disabled",
		"status":  http.StatusOK,
	})
}

// GetWebhookEvents returns events for a webhook
func (h *WebhookHandler) GetWebhookEvents(c *gin.Context) {
	webhookIDStr := c.Param("id")
	webhookID, err := strconv.ParseUint(webhookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid webhook id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageNum, _ := strconv.Atoi(page)
	if pageNum < 1 {
		pageNum = 1
	}

	events, err := h.webhookServ.GetWebhookEvents(c, uint(webhookID), pageNum, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch events",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   events,
		"status": http.StatusOK,
	})
}

// GetAvailableEvents returns all available webhook events
func (h *WebhookHandler) GetAvailableEvents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"events": service.AvailableWebhookEvents(),
		"status": http.StatusOK,
	})
}
