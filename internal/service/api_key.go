package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type APIKeyService struct {
	apiKeyRepo *repository.APIKeyRepository
}

func NewAPIKeyService(apiKeyRepo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}

// CreateAPIKey creates a new API key for a user
func (s *APIKeyService) CreateAPIKey(ctx context.Context, userID uint, name string, scopes []string, expiresAt *time.Time, rateLimit int) (*APIKeyInfo, error) {
	if name == "" {
		return nil, errors.New("API key name is required")
	}

	if rateLimit <= 0 {
		rateLimit = 1000 // Default rate limit
	}

	// Generate API key
	rawKey := generateAPIKey()
	hashedKey := hashAPIKey(rawKey)
	keyPrefix := rawKey[:8]

	// Convert scopes to JSON
	scopesJSON := formatScopes(scopes)

	apiKey := &models.APIKey{
		UserID:    userID,
		Name:      name,
		Key:       hashedKey,
		KeyPrefix: keyPrefix,
		Scopes:    scopesJSON,
		IsActive:  true,
		ExpiresAt: expiresAt,
		RateLimit: rateLimit,
	}

	if err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, fmt.Errorf("error creating API key: %w", err)
	}

	return &APIKeyInfo{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Key:       rawKey, // Only return raw key once
		KeyPrefix: keyPrefix,
		Scopes:    scopes,
		IsActive:  true,
		CreatedAt: apiKey.CreatedAt,
		ExpiresAt: apiKey.ExpiresAt,
		RateLimit: apiKey.RateLimit,
	}, nil
}

type APIKeyInfo struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Key       string     `json:"key,omitempty"` // Only in creation response
	KeyPrefix string     `json:"key_prefix"`
	Scopes    []string   `json:"scopes"`
	IsActive  bool       `json:"is_active"`
	LastUsed  *time.Time `json:"last_used_at"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	RateLimit int        `json:"rate_limit"`
}

// ValidateAPIKey checks if an API key is valid
func (s *APIKeyService) ValidateAPIKey(ctx context.Context, key string) (*models.APIKey, error) {
	hashedKey := hashAPIKey(key)

	apiKey, err := s.apiKeyRepo.GetByKey(ctx, hashedKey)
	if err != nil {
		return nil, errors.New("invalid API key")
	}

	// Check if expired
	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return nil, errors.New("API key expired")
	}

	// Update last used
	_ = s.apiKeyRepo.UpdateLastUsed(ctx, apiKey.ID)

	return apiKey, nil
}

// GetUserAPIKeys returns all API keys for a user (without raw keys)
func (s *APIKeyService) GetUserAPIKeys(ctx context.Context, userID uint) ([]APIKeyInfo, error) {
	keys, err := s.apiKeyRepo.ListUserAPIKeys(ctx, userID)
	if err != nil {
		return nil, err
	}

	infos := make([]APIKeyInfo, len(keys))
	for i, key := range keys {
		scopes := parseScopes(key.Scopes)
		infos[i] = APIKeyInfo{
			ID:        key.ID,
			Name:      key.Name,
			KeyPrefix: key.KeyPrefix,
			Scopes:    scopes,
			IsActive:  key.IsActive,
			LastUsed:  key.LastUsedAt,
			CreatedAt: key.CreatedAt,
			ExpiresAt: key.ExpiresAt,
			RateLimit: key.RateLimit,
		}
	}

	return infos, nil
}

// RevokeAPIKey revokes an API key
func (s *APIKeyService) RevokeAPIKey(ctx context.Context, keyID uint) error {
	return s.apiKeyRepo.Revoke(ctx, keyID)
}

// DeleteAPIKey deletes an API key
func (s *APIKeyService) DeleteAPIKey(ctx context.Context, keyID uint) error {
	return s.apiKeyRepo.Delete(ctx, keyID)
}

// UpdateRateLimit updates the rate limit for an API key
func (s *APIKeyService) UpdateRateLimit(ctx context.Context, keyID uint, rateLimit int) error {
	if rateLimit <= 0 {
		return errors.New("rate limit must be greater than 0")
	}
	return s.apiKeyRepo.UpdateRateLimit(ctx, keyID, rateLimit)
}

// UpdateScopes updates the scopes for an API key
func (s *APIKeyService) UpdateScopes(ctx context.Context, keyID uint, scopes []string) error {
	scopesJSON := formatScopes(scopes)
	return s.apiKeyRepo.UpdateScopes(ctx, keyID, scopesJSON)
}

// CleanupExpiredKeys removes expired API keys
func (s *APIKeyService) CleanupExpiredKeys(ctx context.Context) error {
	return s.apiKeyRepo.DeleteExpiredKeys(ctx)
}

// Helper functions

func generateAPIKey() string {
	// Generate a 32-byte random key and encode as hex (64 chars)
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return "kx_" + hex.EncodeToString(key)
}

func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func formatScopes(scopes []string) string {
	data, _ := json.Marshal(scopes)
	return string(data)
}

func parseScopes(scopesJSON string) []string {
	var scopes []string
	json.Unmarshal([]byte(scopesJSON), &scopes)
	return scopes
}

// AvailableScopes returns all available API scopes
func AvailableScopes() []string {
	return []string{
		"players:read",
		"players:write",
		"bans:read",
		"bans:write",
		"invitations:read",
		"invitations:write",
		"audit:read",
		"sessions:read",
	}
}
