package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type SessionService struct {
	sessionRepo *repository.SessionRepository
	userRepo    *repository.UserRepository
}

func NewSessionService(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

// CreateSession creates a new session for a user
func (s *SessionService) CreateSession(ctx context.Context, userID uint, ipAddress, userAgent string) (*models.Session, error) {
	// Generate session token
	token, err := generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("error generating session token: %w", err)
	}

	// Parse user agent to get device info
	deviceName := parseDeviceName(userAgent)

	session := &models.Session{
		UserID:       userID,
		SessionToken: token,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		DeviceName:   deviceName,
		IsActive:     true,
		LastActivity: time.Now(),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	return session, nil
}

// GetSession retrieves a session by token
func (s *SessionService) GetSession(ctx context.Context, token string) (*models.Session, error) {
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, errors.New("session not found")
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		s.sessionRepo.Logout(ctx, session.ID)
		return nil, errors.New("session expired")
	}

	return session, nil
}

// UpdateSessionActivity updates the last activity timestamp
func (s *SessionService) UpdateSessionActivity(ctx context.Context, sessionID uint) error {
	return s.sessionRepo.UpdateLastActivity(ctx, sessionID)
}

// LogoutSession logs out a session
func (s *SessionService) LogoutSession(ctx context.Context, sessionID uint) error {
	return s.sessionRepo.Logout(ctx, sessionID)
}

// LogoutSessionByToken logs out a session by token
func (s *SessionService) LogoutSessionByToken(ctx context.Context, token string) error {
	return s.sessionRepo.LogoutByToken(ctx, token)
}

// LogoutAllUserSessions logs out all sessions for a user except one
func (s *SessionService) LogoutAllUserSessions(ctx context.Context, userID uint, exceptSessionID uint) error {
	return s.sessionRepo.LogoutAllUserSessions(ctx, userID, exceptSessionID)
}

// LogoutOtherSessions logs out all sessions except the current one
func (s *SessionService) LogoutOtherSessions(ctx context.Context, userID uint, currentSessionToken string) error {
	return s.sessionRepo.LogoutOtherSessions(ctx, userID, currentSessionToken)
}

// GetUserSessions retrieves all sessions for a user
func (s *SessionService) GetUserSessions(ctx context.Context, userID uint) ([]SessionInfo, error) {
	sessions, err := s.sessionRepo.ListUserSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	infos := make([]SessionInfo, len(sessions))
	for i, session := range sessions {
		infos[i] = SessionInfo{
			ID:           session.ID,
			DeviceName:   session.DeviceName,
			IPAddress:    session.IPAddress,
			IsActive:     session.IsActive,
			LastActivity: session.LastActivity,
			LoginAt:      session.LoginAt,
			LogoutAt:     session.LogoutAt,
			ExpiresAt:    session.ExpiresAt,
		}
	}

	return infos, nil
}

// SessionInfo is a public version of Session without token
type SessionInfo struct {
	ID           uint       `json:"id"`
	DeviceName   string     `json:"device_name"`
	IPAddress    string     `json:"ip_address"`
	IsActive     bool       `json:"is_active"`
	LastActivity time.Time  `json:"last_activity"`
	LoginAt      time.Time  `json:"login_at"`
	LogoutAt     *time.Time `json:"logout_at"`
	ExpiresAt    time.Time  `json:"expires_at"`
}

// CleanupExpiredSessions removes expired sessions from database
func (s *SessionService) CleanupExpiredSessions(ctx context.Context) error {
	return s.sessionRepo.DeleteExpiredSessions(ctx)
}

// GetSessionStats returns session statistics for a user
func (s *SessionService) GetSessionStats(ctx context.Context, userID uint) (*SessionStats, error) {
	activeSessions, err := s.sessionRepo.ListActiveUserSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	allSessions, err := s.sessionRepo.ListUserSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Find session with most recent activity (current session)
	var currentSession *models.Session
	if len(activeSessions) > 0 {
		currentSession = &activeSessions[0]
	}

	stats := &SessionStats{
		TotalActiveSessions: int64(len(activeSessions)),
		TotalSessions:       int64(len(allSessions)),
		CurrentDevice:       "",
		CurrentIP:           "",
	}

	if currentSession != nil {
		stats.CurrentDevice = currentSession.DeviceName
		stats.CurrentIP = currentSession.IPAddress
	}

	return stats, nil
}

type SessionStats struct {
	TotalActiveSessions int64  `json:"total_active_sessions"`
	TotalSessions       int64  `json:"total_sessions"`
	CurrentDevice       string `json:"current_device"`
	CurrentIP           string `json:"current_ip"`
}

// Helper functions

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func parseDeviceName(userAgent string) string {
	if userAgent == "" {
		return "Unknown Device"
	}

	// Simple user agent parsing
	browser := "Unknown Browser"
	os := "Unknown OS"

	// Detect browser
	switch {
	case strings.Contains(userAgent, "Chrome"):
		browser = "Chrome"
	case strings.Contains(userAgent, "Firefox"):
		browser = "Firefox"
	case strings.Contains(userAgent, "Safari"):
		browser = "Safari"
	case strings.Contains(userAgent, "Edge"):
		browser = "Edge"
	case strings.Contains(userAgent, "Opera"):
		browser = "Opera"
	}

	// Detect OS
	switch {
	case strings.Contains(userAgent, "Windows"):
		os = "Windows"
	case strings.Contains(userAgent, "Mac"):
		os = "macOS"
	case strings.Contains(userAgent, "Linux"):
		os = "Linux"
	case strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad"):
		os = "iOS"
	case strings.Contains(userAgent, "Android"):
		os = "Android"
	}

	return fmt.Sprintf("%s on %s", browser, os)
}
