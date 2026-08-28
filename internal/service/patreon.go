package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type PatreonService struct {
	patreonRepo *repository.PatreonIntegrationRepository
	playerRepo  *repository.PlayerRepository
	userRepo    *repository.UserRepository
	clientID    string
	clientSecret string
	redirectURI string
}

func NewPatreonService(
	patreonRepo *repository.PatreonIntegrationRepository,
	playerRepo *repository.PlayerRepository,
	userRepo *repository.UserRepository,
) *PatreonService {
	return &PatreonService{
		patreonRepo: patreonRepo,
		playerRepo: playerRepo,
		userRepo: userRepo,
		clientID: os.Getenv("PATREON_CLIENT_ID"),
		clientSecret: os.Getenv("PATREON_CLIENT_SECRET"),
		redirectURI: os.Getenv("PATREON_REDIRECT_URI"),
	}
}

type PatreonOAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type PatreonCampaign struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Name string `json:"name"`
		} `json:"attributes"`
	} `json:"data"`
}

type PatreonMember struct {
	ID         string `json:"id"`
	Attributes struct {
		FullName    string      `json:"full_name"`
		Email       string      `json:"email"`
		IsFollower  bool        `json:"is_follower"`
		PatronStatus string     `json:"patron_status"`
		TierID      string      `json:"currently_entitled_amount_cents"`
	} `json:"attributes"`
	Relationships struct {
		CurrentlyEntitledTier struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		} `json:"currently_entitled_tier"`
	} `json:"relationships"`
}

type PatreonWebhookEvent struct {
	Data struct {
		Type string `json:"type"`
		ID   string `json:"id"`
		Attributes struct {
			Trigger string `json:"trigger"`
		} `json:"attributes"`
	} `json:"data"`
}

type PatreonInfo struct {
	ID               uint                      `json:"id"`
	CampaignID       string                    `json:"campaign_id"`
	IsEnabled        bool                      `json:"is_enabled"`
	LastSyncAt       *time.Time                `json:"last_sync_at"`
	TierMapping      map[string]string         `json:"tier_mapping"`
	MemberCount      int                       `json:"member_count"`
	CreatedAt        time.Time                 `json:"created_at"`
}

// GetOAuthURL generates the Patreon OAuth URL
func (s *PatreonService) GetOAuthURL(state string) string {
	return fmt.Sprintf(
		"https://www.patreon.com/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=identity%%20identity%%5Bemail%%5D%%20campaigns%%20campaigns%%5Bmembers%%5D&state=%s",
		s.clientID,
		url.QueryEscape(s.redirectURI),
		state,
	)
}

// HandleOAuthCallback exchanges the code for an access token
func (s *PatreonService) HandleOAuthCallback(ctx context.Context, userID uint, code string) (*PatreonInfo, error) {
	if s.clientID == "" || s.clientSecret == "" {
		return nil, errors.New("Patreon OAuth not configured")
	}

	// Exchange code for access token
	token, err := s.exchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("error exchanging code: %w", err)
	}

	// Get campaign ID
	campaignID, err := s.getCampaignID(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("error getting campaign: %w", err)
	}

	// Check if integration exists
	existing, _ := s.patreonRepo.GetByUserID(ctx, userID)

	if existing == nil {
		// Create new integration
		integration := &models.PatreonIntegration{
			UserID:       userID,
			CampaignID:   campaignID,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			IsEnabled:    true,
		}
		if err := s.patreonRepo.Create(ctx, integration); err != nil {
			return nil, fmt.Errorf("error creating integration: %w", err)
		}
		return s.getPatreonInfo(integration), nil
	}

	// Update existing integration
	existing.CampaignID = campaignID
	existing.AccessToken = token.AccessToken
	existing.RefreshToken = token.RefreshToken
	existing.IsEnabled = true
	if err := s.patreonRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("error updating integration: %w", err)
	}

	return s.getPatreonInfo(existing), nil
}

// exchangeCodeForToken exchanges OAuth code for access token
func (s *PatreonService) exchangeCodeForToken(code string) (*PatreonOAuthToken, error) {
	data := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"client_id":     {s.clientID},
		"client_secret": {s.clientSecret},
		"redirect_uri":  {s.redirectURI},
	}

	resp, err := http.PostForm("https://www.patreon.com/api/oauth2/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Patreon API error: %s", string(body))
	}

	var token PatreonOAuthToken
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

// getCampaignID retrieves the campaign ID for the authenticated user
func (s *PatreonService) getCampaignID(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://www.patreon.com/api/oauth2/v2/campaigns", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var campaign PatreonCampaign
	if err := json.Unmarshal(body, &campaign); err != nil {
		return "", err
	}

	if len(campaign.Data) == 0 {
		return "", errors.New("no campaigns found")
	}

	return campaign.Data[0].ID, nil
}

// SyncPatreons syncs patrons from Patreon to the system
func (s *PatreonService) SyncPatreons(ctx context.Context, userID uint) error {
	integration, err := s.patreonRepo.GetByUserID(ctx, userID)
	if err != nil || integration == nil {
		return errors.New("Patreon integration not found")
	}

	// Get patrons
	members, err := s.getPatreonMembers(ctx, integration.CampaignID, integration.AccessToken)
	if err != nil {
		return err
	}

	// Get tier mapping
	var tierMap map[string]string
	if integration.RoleMapping != "" {
		if err := json.Unmarshal([]byte(integration.RoleMapping), &tierMap); err != nil {
			tierMap = make(map[string]string)
		}
	} else {
		tierMap = make(map[string]string)
	}

	// Sync each member
	for _, member := range members {
		if member.Attributes.PatronStatus == "active" && member.Relationships.CurrentlyEntitledTier.Data.ID != "" {
			tierID := member.Relationships.CurrentlyEntitledTier.Data.ID
			role, exists := tierMap[tierID]
			if !exists || role == "" {
				role = string(models.ROLE_VIP)
			}

			// Try to find player by VRChat ID (from email or custom field)
			// For now, we'll just log successful sync
			_ = role
		}
	}

	// Update last sync
	now := time.Now()
	integration.LastSyncAt = &now
	if err := s.patreonRepo.Update(ctx, integration); err != nil {
		return err
	}

	return nil
}

// getPatreonMembers retrieves members of a campaign
func (s *PatreonService) getPatreonMembers(ctx context.Context, campaignID, accessToken string) ([]PatreonMember, error) {
	var members []PatreonMember
	cursor := ""

	for {
		url := fmt.Sprintf(
			"https://www.patreon.com/api/oauth2/v2/campaigns/%s/members?include=currently_entitled_tier&fields[member]=full_name,email,patron_status,currently_entitled_amount_cents&page[size]=100",
			campaignID,
		)
		if cursor != "" {
			url += "&page[cursor]=" + cursor
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var result struct {
			Data []PatreonMember `json:"data"`
			Meta struct {
				Pagination struct {
					Cursors struct {
						Next string `json:"next"`
					} `json:"cursors"`
				} `json:"pagination"`
			} `json:"meta"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		members = append(members, result.Data...)

		if result.Meta.Pagination.Cursors.Next == "" {
			break
		}
		cursor = result.Meta.Pagination.Cursors.Next
	}

	return members, nil
}

// ConfigureTierMapping maps Patreon tiers to Kexel roles
func (s *PatreonService) ConfigureTierMapping(ctx context.Context, userID uint, tierMapping map[string]string) error {
	integration, err := s.patreonRepo.GetByUserID(ctx, userID)
	if err != nil || integration == nil {
		return errors.New("Patreon integration not found")
	}

	mappingJSON, err := json.Marshal(tierMapping)
	if err != nil {
		return err
	}

	integration.RoleMapping = string(mappingJSON)
	return s.patreonRepo.Update(ctx, integration)
}

// GetPatreonInfo retrieves Patreon integration info
func (s *PatreonService) GetPatreonInfo(ctx context.Context, userID uint) (*PatreonInfo, error) {
	integration, err := s.patreonRepo.GetByUserID(ctx, userID)
	if err != nil || integration == nil {
		return nil, errors.New("Patreon integration not found")
	}

	return s.getPatreonInfo(integration), nil
}

func (s *PatreonService) getPatreonInfo(integration *models.PatreonIntegration) *PatreonInfo {
	tierMap := make(map[string]string)
	if integration.RoleMapping != "" {
		json.Unmarshal([]byte(integration.RoleMapping), &tierMap)
	}

	return &PatreonInfo{
		ID:          integration.ID,
		CampaignID:  integration.CampaignID,
		IsEnabled:   integration.IsEnabled,
		LastSyncAt:  integration.LastSyncAt,
		TierMapping: tierMap,
		CreatedAt:   integration.CreatedAt,
	}
}

// DisconnectPatreon removes Patreon integration
func (s *PatreonService) DisconnectPatreon(ctx context.Context, userID uint) error {
	integration, err := s.patreonRepo.GetByUserID(ctx, userID)
	if err != nil || integration == nil {
		return errors.New("Patreon integration not found")
	}

	integration.IsEnabled = false
	integration.AccessToken = ""
	integration.RefreshToken = ""
	return s.patreonRepo.Update(ctx, integration)
}

// HandleWebhook processes Patreon webhook events
func (s *PatreonService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	// Verify webhook signature
	// TODO: Implement signature verification using shared secret

	var event PatreonWebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	// Handle different event types
	switch event.Data.Attributes.Trigger {
	case "members:pledge:create":
		// New patron
		return nil
	case "members:pledge:update":
		// Patron tier changed
		return nil
	case "members:pledge:delete":
		// Patron canceled
		return nil
	}

	return nil
}

// GetOAuthConfig returns OAuth configuration URL
func (s *PatreonService) GetOAuthConfig() map[string]string {
	return map[string]string{
		"client_id":    s.clientID,
		"redirect_uri": s.redirectURI,
		"auth_url":     "https://www.patreon.com/oauth2/authorize",
		"token_url":    "https://www.patreon.com/api/oauth2/token",
	}
}
