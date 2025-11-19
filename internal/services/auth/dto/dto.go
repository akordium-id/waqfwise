package dto

import "github.com/akordium-id/waqfwise/internal/shared/domain"

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Name     string      `json:"name"`
	Phone    string      `json:"phone,omitempty"`
	Role     domain.Role `json:"role,omitempty"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	MFACode  string `json:"mfa_code,omitempty"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordConfirmRequest represents reset password confirmation
type ResetPasswordConfirmRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// EnableMFARequest represents MFA enable request
type EnableMFARequest struct {
	Code string `json:"code"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
}

// UserResponse represents user data in response
type UserResponse struct {
	ID         int64       `json:"id"`
	Email      string      `json:"email"`
	Name       string      `json:"name"`
	Phone      string      `json:"phone,omitempty"`
	Role       domain.Role `json:"role"`
	IsActive   bool        `json:"is_active"`
	MFAEnabled bool        `json:"mfa_enabled"`
	TenantID   *int64      `json:"tenant_id,omitempty"`
}

// MFASetupResponse represents MFA setup response
type MFASetupResponse struct {
	Secret     string `json:"secret"`
	QRCodeURL  string `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// FromDomain converts domain.User to UserResponse
func FromDomain(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Phone:      user.Phone,
		Role:       user.Role,
		IsActive:   user.IsActive,
		MFAEnabled: user.MFAEnabled,
		TenantID:   user.TenantID,
	}
}
