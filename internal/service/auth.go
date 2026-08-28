package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo          *repository.UserRepository
	refreshTokenRepo  *repository.RefreshTokenRepository
	invitationRepo    *repository.InvitationRepository
	jwtSecret         string
	accessTokenExpiry time.Duration
	refreshTokenExp   time.Duration
}

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		refreshTokenRepo:  refreshTokenRepo,
		invitationRepo:    nil,
		jwtSecret:         jwtSecret,
		accessTokenExpiry: 15 * time.Minute,
		refreshTokenExp:   7 * 24 * time.Hour,
	}
}

// NewAuthServiceWithInvitations creates auth service with invitation support
func NewAuthServiceWithInvitations(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	invitationRepo *repository.InvitationRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		refreshTokenRepo:  refreshTokenRepo,
		invitationRepo:    invitationRepo,
		jwtSecret:         jwtSecret,
		accessTokenExpiry: 15 * time.Minute,
		refreshTokenExp:   7 * 24 * time.Hour,
	}
}

// Register creates a new user
func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingEmail, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     models.ROLE_USER,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

// Login authenticates user and returns tokens and user
func (s *AuthService) Login(ctx context.Context, username, password string) (accessToken, refreshToken string, user *models.User, err error) {
	user, err = s.userRepo.GetByUsername(ctx, username)
	if err != nil || user == nil {
		return "", "", nil, errors.New("invalid username or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", nil, errors.New("invalid username or password")
	}

	// Generate access token
	accessToken, err = s.generateAccessToken(user)
	if err != nil {
		return "", "", nil, err
	}

	// Generate and save refresh token
	refreshToken, err = s.generateAndSaveRefreshToken(ctx, user)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

// RefreshAccessToken generates new access token from refresh token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenStr string) (accessToken string, err error) {
	// Verify refresh token in database
	refreshToken, err := s.refreshTokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil || refreshToken == nil {
		return "", errors.New("invalid refresh token")
	}

	// Check if token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		s.refreshTokenRepo.Delete(ctx, refreshToken.ID)
		return "", errors.New("refresh token expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, refreshToken.UserID)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	// Generate new access token
	accessToken, err = s.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// VerifyAccessToken verifies and parses access token
func (s *AuthService) VerifyAccessToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) generateAndSaveRefreshToken(ctx context.Context, user *models.User) (string, error) {
	// Generate random token
	randomToken := jwt.New(jwt.SigningMethodHS256)
	randomClaims := randomToken.Claims.(jwt.MapClaims)
	randomClaims["exp"] = time.Now().Add(s.refreshTokenExp).Unix()
	randomClaims["user_id"] = user.ID

	tokenString, err := randomToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	// Save to database
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(s.refreshTokenExp),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return "", err
	}

	return tokenString, nil
}

// Logout invalidates refresh token
func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string) error {
	refreshToken, err := s.refreshTokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil || refreshToken == nil {
		return errors.New("invalid refresh token")
	}

	return s.refreshTokenRepo.Delete(ctx, refreshToken.ID)
}
