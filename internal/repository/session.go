package repository

import (
	"context"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByToken retrieves a session by token
func (r *SessionRepository) GetByToken(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	err := r.db.WithContext(ctx).Where("session_token = ? AND is_active = true", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByID retrieves a session by ID
func (r *SessionRepository) GetByID(ctx context.Context, id uint) (*models.Session, error) {
	var session models.Session
	err := r.db.WithContext(ctx).First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// ListUserSessions retrieves all sessions for a user
func (r *SessionRepository) ListUserSessions(ctx context.Context, userID uint) ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("last_activity DESC").
		Find(&sessions).Error
	return sessions, err
}

// ListActiveUserSessions retrieves all active sessions for a user
func (r *SessionRepository) ListActiveUserSessions(ctx context.Context, userID uint) ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = true", userID).
		Order("last_activity DESC").
		Find(&sessions).Error
	return sessions, err
}

// UpdateLastActivity updates the last activity timestamp
func (r *SessionRepository) UpdateLastActivity(ctx context.Context, sessionID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("last_activity", time.Now()).Error
}

// Logout marks a session as inactive
func (r *SessionRepository) Logout(ctx context.Context, sessionID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"is_active": false,
			"logout_at": now,
		}).Error
}

// LogoutByToken marks a session as inactive by token
func (r *SessionRepository) LogoutByToken(ctx context.Context, token string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("session_token = ?", token).
		Updates(map[string]interface{}{
			"is_active": false,
			"logout_at": now,
		}).Error
}

// LogoutAllUserSessions logs out all sessions for a user (except current)
func (r *SessionRepository) LogoutAllUserSessions(ctx context.Context, userID uint, excludeSessionID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ? AND id != ?", userID, excludeSessionID).
		Updates(map[string]interface{}{
			"is_active": false,
			"logout_at": now,
		}).Error
}

// LogoutOtherSessions logs out all sessions for a user except the provided one
func (r *SessionRepository) LogoutOtherSessions(ctx context.Context, userID uint, currentSessionToken string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ? AND session_token != ?", userID, currentSessionToken).
		Updates(map[string]interface{}{
			"is_active": false,
			"logout_at": now,
		}).Error
}

// DeleteExpiredSessions deletes sessions that have expired
func (r *SessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.Session{}).Error
}

// GetSessionCount returns the count of active sessions for a user
func (r *SessionRepository) GetSessionCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ? AND is_active = true", userID).
		Count(&count).Error
	return count, err
}

// CheckSessionValid checks if a session is still valid
func (r *SessionRepository) CheckSessionValid(ctx context.Context, token string) (bool, error) {
	var session models.Session
	err := r.db.WithContext(ctx).
		Where("session_token = ? AND is_active = true AND expires_at > ?", token, time.Now()).
		First(&session).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}
