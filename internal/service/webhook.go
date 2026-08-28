package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type WebhookService struct {
	webhookRepo      *repository.WebhookRepository
	webhookEventRepo *repository.WebhookEventRepository
	httpClient       *http.Client
}

func NewWebhookService(webhookRepo *repository.WebhookRepository, webhookEventRepo *repository.WebhookEventRepository) *WebhookService {
	return &WebhookService{
		webhookRepo:      webhookRepo,
		webhookEventRepo: webhookEventRepo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateWebhook creates a new webhook
func (s *WebhookService) CreateWebhook(ctx context.Context, userID uint, name string, url string, events []string) (*WebhookInfo, error) {
	if name == "" || url == "" || len(events) == 0 {
		return nil, errors.New("missing required fields")
	}

	secret := generateWebhookSecret()
	eventsJSON, _ := json.Marshal(events)

	webhook := &models.Webhook{
		UserID: userID,
		Name:   name,
		URL:    url,
		Events: string(eventsJSON),
		Secret: secret,
		IsActive: true,
	}

	if err := s.webhookRepo.Create(ctx, webhook); err != nil {
		return nil, fmt.Errorf("error creating webhook: %w", err)
	}

	return &WebhookInfo{
		ID:        webhook.ID,
		Name:      webhook.Name,
		URL:       webhook.URL,
		Events:    events,
		IsActive:  webhook.IsActive,
		CreatedAt: webhook.CreatedAt,
	}, nil
}

type WebhookInfo struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	Events        []string   `json:"events"`
	IsActive      bool       `json:"is_active"`
	FailureCount  int        `json:"failure_count"`
	LastTriedAt   *time.Time `json:"last_tried_at"`
	LastSuccessAt *time.Time `json:"last_success_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

// GetUserWebhooks returns all webhooks for a user
func (s *WebhookService) GetUserWebhooks(ctx context.Context, userID uint) ([]WebhookInfo, error) {
	webhooks, err := s.webhookRepo.ListUserWebhooks(ctx, userID)
	if err != nil {
		return nil, err
	}

	infos := make([]WebhookInfo, len(webhooks))
	for i, wh := range webhooks {
		var events []string
		json.Unmarshal([]byte(wh.Events), &events)
		infos[i] = WebhookInfo{
			ID:            wh.ID,
			Name:          wh.Name,
			URL:           wh.URL,
			Events:        events,
			IsActive:      wh.IsActive,
			FailureCount:  wh.FailureCount,
			LastTriedAt:   wh.LastTriedAt,
			LastSuccessAt: wh.LastSuccessAt,
			CreatedAt:     wh.CreatedAt,
		}
	}

	return infos, nil
}

// DeleteWebhook deletes a webhook
func (s *WebhookService) DeleteWebhook(ctx context.Context, id uint) error {
	return s.webhookRepo.Delete(ctx, id)
}

// DisableWebhook disables a webhook
func (s *WebhookService) DisableWebhook(ctx context.Context, id uint) error {
	webhook, err := s.webhookRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	webhook.IsActive = false
	return s.webhookRepo.Update(ctx, webhook)
}

// TriggerWebhookEvent triggers a webhook event
func (s *WebhookService) TriggerWebhookEvent(ctx context.Context, eventType string, payload interface{}) error {
	webhooks, err := s.webhookRepo.ListActiveWebhooks(ctx)
	if err != nil {
		return err
	}

	payloadJSON, _ := json.Marshal(payload)

	for _, webhook := range webhooks {
		var events []string
		json.Unmarshal([]byte(webhook.Events), &events)

		// Check if webhook is subscribed to this event
		subscribed := false
		for _, e := range events {
			if e == eventType || e == "*" {
				subscribed = true
				break
			}
		}

		if !subscribed {
			continue
		}

		// Create webhook event record
		event := &models.WebhookEvent{
			WebhookID: webhook.ID,
			EventType: eventType,
			Payload:   string(payloadJSON),
			Attempts:  0,
		}

		if err := s.webhookEventRepo.Create(ctx, event); err != nil {
			continue
		}

		// Attempt delivery
		go s.deliverWebhookEvent(&webhook, event, payloadJSON)
	}

	return nil
}

// deliverWebhookEvent attempts to deliver a webhook event
func (s *WebhookService) deliverWebhookEvent(webhook *models.Webhook, event *models.WebhookEvent, payload []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	signature := generateSignature(webhook.Secret, payload)

	req, _ := http.NewRequestWithContext(ctx, "POST", webhook.URL, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-Event", event.EventType)

	s.webhookRepo.UpdateLastTried(ctx, webhook.ID)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		event.StatusCode = 0
		event.Response = err.Error()
		event.Attempts++
		event.NextRetry = timePtr(time.Now().Add(time.Duration(event.Attempts*5) * time.Minute))
		s.webhookEventRepo.UpdateEvent(ctx, event)

		s.webhookRepo.UpdateFailureCount(ctx, webhook.ID, webhook.FailureCount+1)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	event.StatusCode = resp.StatusCode
	event.Response = string(body)
	event.Attempts++

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		event.IsDelivered = true
		s.webhookEventRepo.MarkAsDelivered(ctx, event.ID)
		s.webhookRepo.UpdateLastSuccess(ctx, webhook.ID)
	} else {
		// Retry later
		event.NextRetry = timePtr(time.Now().Add(time.Duration(event.Attempts*5) * time.Minute))
		s.webhookRepo.UpdateFailureCount(ctx, webhook.ID, webhook.FailureCount+1)
	}

	s.webhookEventRepo.UpdateEvent(ctx, event)
}

// GetWebhookEvents returns paginated events for a webhook
func (s *WebhookService) GetWebhookEvents(ctx context.Context, webhookID uint, page int, pageSize int) ([]WebhookEventInfo, error) {
	offset := (page - 1) * pageSize
	events, err := s.webhookEventRepo.ListWebhookEvents(ctx, webhookID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	infos := make([]WebhookEventInfo, len(events))
	for i, e := range events {
		infos[i] = WebhookEventInfo{
			ID:          e.ID,
			EventType:   e.EventType,
			StatusCode:  e.StatusCode,
			IsDelivered: e.IsDelivered,
			Attempts:    e.Attempts,
			CreatedAt:   e.CreatedAt,
		}
	}

	return infos, nil
}

type WebhookEventInfo struct {
	ID          uint       `json:"id"`
	EventType   string     `json:"event_type"`
	StatusCode  int        `json:"status_code"`
	IsDelivered bool       `json:"is_delivered"`
	Attempts    int        `json:"attempts"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Helper functions

func generateWebhookSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateSignature(secret string, payload []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// AvailableWebhookEvents returns all available webhook events
func AvailableWebhookEvents() []string {
	return []string{
		"player.created",
		"player.updated",
		"player.deleted",
		"ban.created",
		"ban.updated",
		"ban.deleted",
		"invitation.created",
		"invitation.used",
		"invitation.revoked",
		"*", // All events
	}
}
