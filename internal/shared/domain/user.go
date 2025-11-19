package domain

import (
	"time"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleNazir    Role = "nazir"
	RoleDonor    Role = "donor"
	RoleAuditor  Role = "auditor"
	RoleOperator Role = "operator"
)

// User represents a user in the system
type User struct {
	ID           int64     `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name"`
	Phone        string    `json:"phone,omitempty" db:"phone"`
	Role         Role      `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	MFAEnabled   bool      `json:"mfa_enabled" db:"mfa_enabled"`
	MFASecret    string    `json:"-" db:"mfa_secret"`
	TenantID     *int64    `json:"tenant_id,omitempty" db:"tenant_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

// UserProfile represents additional user profile information
type UserProfile struct {
	UserID      int64     `json:"user_id" db:"user_id"`
	Avatar      string    `json:"avatar,omitempty" db:"avatar"`
	Address     string    `json:"address,omitempty" db:"address"`
	City        string    `json:"city,omitempty" db:"city"`
	Country     string    `json:"country,omitempty" db:"country"`
	PostalCode  string    `json:"postal_code,omitempty" db:"postal_code"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role Role) bool {
	return u.Role == role
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// CanAccessEnterprise checks if user can access enterprise features
func (u *User) CanAccessEnterprise() bool {
	return u.TenantID != nil
}
