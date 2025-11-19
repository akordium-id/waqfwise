package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/akordium-id/waqfwise/internal/services/auth/dto"
	"github.com/akordium-id/waqfwise/internal/services/auth/repository"
	"github.com/akordium-id/waqfwise/internal/shared/domain"
	"github.com/akordium-id/waqfwise/internal/shared/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// Service defines auth service interface
type Service interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	GetProfile(ctx context.Context, userID int64) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, userID int64, req *dto.ChangePasswordRequest) error
	SetupMFA(ctx context.Context, userID int64) (*dto.MFASetupResponse, error)
	EnableMFA(ctx context.Context, userID int64, code string) error
	DisableMFA(ctx context.Context, userID int64) error
	ValidateToken(token string) (*Claims, error)
}

type service struct {
	repo      repository.Repository
	jwtSecret []byte
}

// Claims represents JWT claims
type Claims struct {
	UserID   int64       `json:"user_id"`
	Email    string      `json:"email"`
	Role     domain.Role `json:"role"`
	TenantID *int64      `json:"tenant_id,omitempty"`
	jwt.StandardClaims
}

// New creates a new auth service
func New(repo repository.Repository, jwtSecret string) Service {
	return &service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register registers a new user
func (s *service) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, errors.New(errors.ErrCodeValidation, "Email, password, and name are required", 400)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "Failed to hash password", 500)
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = domain.RoleDonor
	}

	// Create user
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Phone:        req.Phone,
		Role:         role,
		IsActive:     true,
		MFAEnabled:   false,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:         dto.FromDomain(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// Login authenticates a user
func (s *service) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New(errors.ErrCodeForbidden, "Account is inactive", 403)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Verify MFA if enabled
	if user.MFAEnabled {
		if req.MFACode == "" {
			return nil, errors.New(errors.ErrCodeUnauthorized, "MFA code is required", 401)
		}

		if !totp.Validate(req.MFACode, user.MFASecret) {
			return nil, errors.New(errors.ErrCodeUnauthorized, "Invalid MFA code", 401)
		}
	}

	// Update last login
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:         dto.FromDomain(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// RefreshToken refreshes access token
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// Get user
	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New(errors.ErrCodeForbidden, "Account is inactive", 403)
	}

	// Generate new tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:         dto.FromDomain(user),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// GetProfile gets user profile
func (s *service) GetProfile(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.FromDomain(user), nil
}

// ChangePassword changes user password
func (s *service) ChangePassword(ctx context.Context, userID int64, req *dto.ChangePasswordRequest) error {
	// Get user
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New(errors.ErrCodeBadRequest, "Invalid old password", 400)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to hash password", 500)
	}

	// Update password
	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

// SetupMFA sets up MFA for user
func (s *service) SetupMFA(ctx context.Context, userID int64) (*dto.MFASetupResponse, error) {
	// Get user
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Generate secret
	secret := make([]byte, 20)
	_, err = rand.Read(secret)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "Failed to generate secret", 500)
	}

	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	// Save secret
	if err := s.repo.SetMFASecret(ctx, userID, secretBase32); err != nil {
		return nil, err
	}

	// Generate QR code URL
	qrCodeURL := fmt.Sprintf(
		"otpauth://totp/WaqfWise:%s?secret=%s&issuer=WaqfWise",
		user.Email,
		secretBase32,
	)

	// Generate backup codes
	backupCodes := make([]string, 10)
	for i := 0; i < 10; i++ {
		code := make([]byte, 4)
		rand.Read(code)
		backupCodes[i] = base32.StdEncoding.EncodeToString(code)
	}

	return &dto.MFASetupResponse{
		Secret:     secretBase32,
		QRCodeURL:  qrCodeURL,
		BackupCodes: backupCodes,
	}, nil
}

// EnableMFA enables MFA for user
func (s *service) EnableMFA(ctx context.Context, userID int64, code string) error {
	// Get user
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.MFASecret == "" {
		return errors.New(errors.ErrCodeBadRequest, "MFA not set up", 400)
	}

	// Verify code
	if !totp.Validate(code, user.MFASecret) {
		return errors.New(errors.ErrCodeBadRequest, "Invalid MFA code", 400)
	}

	// Enable MFA
	return s.repo.EnableMFA(ctx, userID)
}

// DisableMFA disables MFA for user
func (s *service) DisableMFA(ctx context.Context, userID int64) error {
	return s.repo.DisableMFA(ctx, userID)
}

// ValidateToken validates JWT token
func (s *service) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInvalidToken, "Invalid token", 401)
	}

	if !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}

// generateAccessToken generates access token
func (s *service) generateAccessToken(user *domain.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
		TenantID: user.TenantID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "waqfwise",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// generateRefreshToken generates refresh token
func (s *service) generateRefreshToken(user *domain.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
		TenantID: user.TenantID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
			IssuedAt:  time.Now().Unix(),
			Issuer:    "waqfwise",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
