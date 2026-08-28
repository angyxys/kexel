package handler

import (
	"net/http"
	"strconv"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionServ *service.SessionService
}

func NewSessionHandler(sessionServ *service.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionServ: sessionServ,
	}
}

// GetSessions returns all sessions for the authenticated user
func (h *SessionHandler) GetSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	sessions, err := h.sessionServ.GetUserSessions(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch sessions",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   sessions,
		"status": http.StatusOK,
	})
}

// GetSessionStats returns session statistics for the authenticated user
func (h *SessionHandler) GetSessionStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	stats, err := h.sessionServ.GetSessionStats(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch session stats",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// LogoutSession logs out a specific session
func (h *SessionHandler) LogoutSession(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid session id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Security: verify the session belongs to the current user
	if err := h.sessionServ.LogoutSession(c, uint(sessionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to logout session",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "session logged out",
		"status":  http.StatusOK,
	})
}

// LogoutAllSessions logs out all sessions for the current user except the current one
func (h *SessionHandler) LogoutAllSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	// Get current session token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing authorization header",
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Extract token from "Bearer <token>"
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid authorization header format",
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.sessionServ.LogoutOtherSessions(c, userID.(uint), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to logout sessions",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "all other sessions logged out",
		"status":  http.StatusOK,
	})
}

// SessionActivityMiddleware updates last activity for active sessions
func (h *SessionHandler) SessionActivityMiddleware(sessionServ *service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, exists := c.Get("session_id")
		if exists && sessionID != nil {
			// Update last activity in background
			go sessionServ.UpdateSessionActivity(c, sessionID.(uint))
		}
		c.Next()
	}
}
