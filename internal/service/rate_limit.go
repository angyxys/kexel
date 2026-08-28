package service

import (
	"context"
	"sync"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type RateLimitService struct {
	ruleRepo *repository.RateLimitRuleRepository
	logRepo  *repository.RateLimitLogRepository
	// In-memory cache for performance
	requestCounts map[string]*RequestCounter
	mu             sync.RWMutex
}

type RequestCounter struct {
	Count     int
	WindowEnd time.Time
}

func NewRateLimitService(ruleRepo *repository.RateLimitRuleRepository, logRepo *repository.RateLimitLogRepository) *RateLimitService {
	svc := &RateLimitService{
		ruleRepo:      ruleRepo,
		logRepo:       logRepo,
		requestCounts: make(map[string]*RequestCounter),
	}

	// Cleanup expired counters every minute
	go svc.cleanupExpiredCounters()

	return svc
}

// CheckRateLimit checks if a request should be allowed
func (s *RateLimitService) CheckRateLimit(ctx context.Context, ip string, endpoint string) (bool, error) {
	rule, err := s.ruleRepo.GetByIP(ctx, ip)
	if err != nil {
		// No rule = no limit
		return true, nil
	}

	key := ip + ":" + endpoint
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	counter, exists := s.requestCounts[key]

	// Check if window has expired
	if !exists || now.After(counter.WindowEnd) {
		s.requestCounts[key] = &RequestCounter{
			Count:     1,
			WindowEnd: now.Add(time.Duration(rule.WindowSize) * time.Second),
		}
		return true, nil
	}

	// Increment counter
	counter.Count++

	// Check if limit exceeded
	if counter.Count > rule.RequestLimit {
		// Log the rate limit hit
		_ = s.logRepo.Create(ctx, &models.RateLimitLog{
			RuleID:       rule.ID,
			IPAddress:    ip,
			Endpoint:     endpoint,
			RequestCount: counter.Count,
		})
		return false, nil
	}

	return true, nil
}

// CreateRule creates a new rate limit rule
func (s *RateLimitService) CreateRule(ctx context.Context, ip string, endpoint string, limit int, window int) (*RateLimitRule, error) {
	rule := &models.RateLimitRule{
		IPAddress:    ip,
		Endpoint:     endpoint,
		RequestLimit: limit,
		WindowSize:   window,
		IsActive:     true,
	}

	if err := s.ruleRepo.Create(ctx, rule); err != nil {
		return nil, err
	}

	return &RateLimitRule{
		ID:           rule.ID,
		IPAddress:    rule.IPAddress,
		RequestLimit: rule.RequestLimit,
		WindowSize:   rule.WindowSize,
	}, nil
}

type RateLimitRule struct {
	ID           uint `json:"id"`
	IPAddress    string `json:"ip_address"`
	RequestLimit int `json:"request_limit"`
	WindowSize   int `json:"window_size"`
}

// DeleteRule deletes a rate limit rule
func (s *RateLimitService) DeleteRule(ctx context.Context, ip string) error {
	return s.ruleRepo.Delete(ctx, ip)
}

// GetRecentBlocks gets recent rate limit blocks
func (s *RateLimitService) GetRecentBlocks(ctx context.Context, limit int) ([]RateLimitBlock, error) {
	logs, err := s.logRepo.GetRecentBlocks(ctx, limit)
	if err != nil {
		return nil, err
	}

	blocks := make([]RateLimitBlock, len(logs))
	for i, log := range logs {
		blocks[i] = RateLimitBlock{
			IPAddress:    log.IPAddress,
			Endpoint:     log.Endpoint,
			RequestCount: log.RequestCount,
			BlockedAt:    log.BlockedAt,
		}
	}

	return blocks, nil
}

type RateLimitBlock struct {
	IPAddress    string    `json:"ip_address"`
	Endpoint     string    `json:"endpoint"`
	RequestCount int       `json:"request_count"`
	BlockedAt    time.Time `json:"blocked_at"`
}

// cleanupExpiredCounters removes expired counters from memory
func (s *RateLimitService) cleanupExpiredCounters() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, counter := range s.requestCounts {
			if now.After(counter.WindowEnd) {
				delete(s.requestCounts, key)
			}
		}
		s.mu.Unlock()
	}
}
