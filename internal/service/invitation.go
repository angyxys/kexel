package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type InvitationService struct {
	repo *repository.InvitationRepository
}

func NewInvitationService(repo *repository.InvitationRepository) *InvitationService {
	return &InvitationService{repo: repo}
}

// GenerateCode creates a random invitation code
func (s *InvitationService) GenerateCode() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateInvitation creates a new invitation code
type InvitationRequest struct {
	Role      string
	MaxUses   int
	ExpiresAt *time.Time
}

func (s *InvitationService) CreateInvitation(ctx context.Context, createdBy uint, req InvitationRequest) (*models.InvitationCode, error) {
	code, err := s.GenerateCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}

	invitation := &models.InvitationCode{
		Code:      code,
		CreatedBy: createdBy,
		Role:      req.Role,
		MaxUses:   req.MaxUses,
		ExpiresAt: req.ExpiresAt,
		IsActive:  true,
	}

	if err := s.repo.Create(ctx, invitation); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	return invitation, nil
}

// ValidateInvitation checks if an invitation code is valid
func (s *InvitationService) ValidateInvitation(ctx context.Context, code string) (*models.InvitationCode, error) {
	invitation, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if invitation == nil {
		return nil, fmt.Errorf("invitation code not found")
	}

	if !invitation.IsActive {
		return nil, fmt.Errorf("invitation code is inactive")
	}

	if invitation.ExpiresAt != nil && time.Now().After(*invitation.ExpiresAt) {
		return nil, fmt.Errorf("invitation code has expired")
	}

	if invitation.MaxUses != -1 && invitation.Uses >= invitation.MaxUses {
		return nil, fmt.Errorf("invitation code usage limit reached")
	}

	return invitation, nil
}

// UseInvitation records usage of an invitation code
func (s *InvitationService) UseInvitation(ctx context.Context, code string) error {
	invitation, err := s.ValidateInvitation(ctx, code)
	if err != nil {
		return err
	}

	if err := s.repo.IncrementUses(ctx, invitation.ID); err != nil {
		return fmt.Errorf("failed to record usage: %w", err)
	}

	return nil
}

// RevokeInvitation deactivates an invitation code
func (s *InvitationService) RevokeInvitation(ctx context.Context, id uint) error {
	return s.repo.Revoke(ctx, id)
}

// GetInvitations retrieves invitations created by a user
func (s *InvitationService) GetInvitations(ctx context.Context, createdBy uint) ([]models.InvitationCode, error) {
	return s.repo.ListByCreator(ctx, createdBy)
}

// GetActiveInvitations retrieves all active invitations
func (s *InvitationService) GetActiveInvitations(ctx context.Context) ([]models.InvitationCode, error) {
	return s.repo.ListActive(ctx)
}

// InvitationInfo contains information about an invitation
type InvitationInfo struct {
	ID        uint      `json:"id"`
	Code      string    `json:"code"`
	Role      string    `json:"role"`
	MaxUses   int       `json:"max_uses"`
	Uses      int       `json:"uses"`
	ExpiresAt *time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
	IsValid   bool      `json:"is_valid"`
	CreatedAt time.Time `json:"created_at"`
	Creator   string    `json:"creator"`
}

// GetInvitationInfo gets detailed info about an invitation
func (s *InvitationService) GetInvitationInfo(ctx context.Context, id uint) (*InvitationInfo, error) {
	invitation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if invitation == nil {
		return nil, fmt.Errorf("invitation not found")
	}

	isValid, _ := s.repo.IsValid(ctx, invitation.Code)

	return &InvitationInfo{
		ID:        invitation.ID,
		Code:      invitation.Code,
		Role:      invitation.Role,
		MaxUses:   invitation.MaxUses,
		Uses:      invitation.Uses,
		ExpiresAt: invitation.ExpiresAt,
		IsActive:  invitation.IsActive,
		IsValid:   isValid,
		CreatedAt: invitation.CreatedAt,
		Creator:   invitation.Creator.Username,
	}, nil
}
