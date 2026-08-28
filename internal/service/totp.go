package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
	"github.com/pquerna/otp/totp"
)

type TOTPService struct {
	totpRepo *repository.TOTPRepository
}

func NewTOTPService(totpRepo *repository.TOTPRepository) *TOTPService {
	return &TOTPService{
		totpRepo: totpRepo,
	}
}

// GenerateTOTPSecret generates a new TOTP secret for a user
func (s *TOTPService) GenerateTOTPSecret(ctx context.Context, userID uint, email string) (*TOTPSetup, error) {
	// Generate new TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Kexel",
		AccountName: email,
	})
	if err != nil {
		return nil, fmt.Errorf("error generating TOTP key: %w", err)
	}

	// Generate backup codes
	backupCodes, err := generateBackupCodes(10)
	if err != nil {
		return nil, fmt.Errorf("error generating backup codes: %w", err)
	}

	// Store in database (not enabled yet)
	backupCodesJSON, _ := json.Marshal(backupCodes)
	totpSecret := &models.TOTPSecret{
		UserID:      userID,
		Secret:      key.Secret(),
		BackupCodes: string(backupCodesJSON),
		IsEnabled:   false,
	}

	if err := s.totpRepo.Create(ctx, totpSecret); err != nil {
		return nil, fmt.Errorf("error saving TOTP secret: %w", err)
	}

	return &TOTPSetup{
		Secret:      key.Secret(),
		QRCode:      key.URL(),
		BackupCodes: backupCodes,
	}, nil
}

// TOTPSetup contains the setup information for 2FA
type TOTPSetup struct {
	Secret      string   `json:"secret"`
	QRCode      string   `json:"qr_code"`
	BackupCodes []string `json:"backup_codes"`
}

// EnableTOTP enables TOTP for a user after verification
func (s *TOTPService) EnableTOTP(ctx context.Context, userID uint, code string) error {
	totpSecret, err := s.totpRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("TOTP secret not found")
	}

	if totpSecret.IsEnabled {
		return fmt.Errorf("TOTP is already enabled")
	}

	// Verify the code
	if !totp.Validate(code, totpSecret.Secret) {
		return fmt.Errorf("invalid code")
	}

	// Enable TOTP
	now := time.Now()
	return s.totpRepo.EnableTOTP(ctx, userID, &now)
}

// VerifyTOTPCode verifies a TOTP code
func (s *TOTPService) VerifyTOTPCode(ctx context.Context, userID uint, code string) (bool, error) {
	totpSecret, err := s.totpRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("TOTP secret not found")
	}

	if !totpSecret.IsEnabled {
		return false, fmt.Errorf("TOTP is not enabled")
	}

	// Verify the code with a window of ±1 (30-second intervals)
	return totp.Validate(code, totpSecret.Secret), nil
}

// VerifyBackupCode verifies and consumes a backup code
func (s *TOTPService) VerifyBackupCode(ctx context.Context, userID uint, code string) (bool, error) {
	totpSecret, err := s.totpRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("TOTP secret not found")
	}

	if !totpSecret.IsEnabled {
		return false, fmt.Errorf("TOTP is not enabled")
	}

	var backupCodes []string
	if err := json.Unmarshal([]byte(totpSecret.BackupCodes), &backupCodes); err != nil {
		return false, fmt.Errorf("error parsing backup codes")
	}

	// Check if code exists and remove it
	found := false
	for i, bc := range backupCodes {
		if bc == code {
			// Remove used code
			backupCodes = append(backupCodes[:i], backupCodes[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return false, nil
	}

	// Update backup codes
	updatedCodes, _ := json.Marshal(backupCodes)
	return true, s.totpRepo.UpdateBackupCodes(ctx, userID, string(updatedCodes))
}

// DisableTOTP disables TOTP for a user
func (s *TOTPService) DisableTOTP(ctx context.Context, userID uint) error {
	return s.totpRepo.Delete(ctx, userID)
}

// IsTOTPEnabled checks if TOTP is enabled for a user
func (s *TOTPService) IsTOTPEnabled(ctx context.Context, userID uint) (bool, error) {
	totpSecret, err := s.totpRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, nil // Not found means not enabled
	}
	return totpSecret.IsEnabled, nil
}

// GetTOTPStatus returns TOTP status for a user
func (s *TOTPService) GetTOTPStatus(ctx context.Context, userID uint) (*TOTPStatus, error) {
	totpSecret, err := s.totpRepo.GetByUserID(ctx, userID)
	if err != nil {
		return &TOTPStatus{
			IsEnabled:     false,
			BackupCodesLeft: 0,
		}, nil
	}

	var backupCodes []string
	json.Unmarshal([]byte(totpSecret.BackupCodes), &backupCodes)

	return &TOTPStatus{
		IsEnabled:       totpSecret.IsEnabled,
		BackupCodesLeft: int64(len(backupCodes)),
		EnabledAt:       totpSecret.EnabledAt,
		LastUsedAt:      totpSecret.LastUsedAt,
	}, nil
}

type TOTPStatus struct {
	IsEnabled       bool       `json:"is_enabled"`
	BackupCodesLeft int64      `json:"backup_codes_left"`
	EnabledAt       *time.Time `json:"enabled_at"`
	LastUsedAt      *time.Time `json:"last_used_at"`
}

// Helper functions

func generateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		code, err := generateRandomCode(8)
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}
	return codes, nil
}

func generateRandomCode(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}
