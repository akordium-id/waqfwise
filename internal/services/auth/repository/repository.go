package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/akordium-id/waqfwise/internal/shared/domain"
	"github.com/akordium-id/waqfwise/internal/shared/errors"
)

// Repository defines auth repository interface
type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, userID int64, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userID int64) error
	SetMFASecret(ctx context.Context, userID int64, secret string) error
	EnableMFA(ctx context.Context, userID int64) error
	DisableMFA(ctx context.Context, userID int64) error
}

type repository struct {
	db *sql.DB
}

// New creates a new auth repository
func New(db *sql.DB) Repository {
	return &repository{db: db}
}

// Create creates a new user
func (r *repository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, name, phone, role, is_active, mfa_enabled, tenant_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Phone,
		user.Role,
		user.IsActive,
		user.MFAEnabled,
		user.TenantID,
		now,
		now,
	).Scan(&user.ID)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return errors.Wrap(err, errors.ErrCodeDuplicateEntry, "Email already exists", 409)
		}
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to create user", 500)
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// FindByID finds a user by ID
func (r *repository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, phone, role, is_active, mfa_enabled, mfa_secret,
		       tenant_id, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	var lastLoginAt sql.NullTime
	var phone, mfaSecret sql.NullString
	var tenantID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&phone,
		&user.Role,
		&user.IsActive,
		&user.MFAEnabled,
		&mfaSecret,
		&tenantID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrCodeNotFound, "User not found", 404)
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "Failed to find user", 500)
	}

	if phone.Valid {
		user.Phone = phone.String
	}
	if mfaSecret.Valid {
		user.MFASecret = mfaSecret.String
	}
	if tenantID.Valid {
		id := tenantID.Int64
		user.TenantID = &id
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByEmail finds a user by email
func (r *repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, phone, role, is_active, mfa_enabled, mfa_secret,
		       tenant_id, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	var lastLoginAt sql.NullTime
	var phone, mfaSecret sql.NullString
	var tenantID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&phone,
		&user.Role,
		&user.IsActive,
		&user.MFAEnabled,
		&mfaSecret,
		&tenantID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrCodeNotFound, "User not found", 404)
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "Failed to find user", 500)
	}

	if phone.Valid {
		user.Phone = phone.String
	}
	if mfaSecret.Valid {
		user.MFASecret = mfaSecret.String
	}
	if tenantID.Valid {
		id := tenantID.Int64
		user.TenantID = &id
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// Update updates a user
func (r *repository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, phone = $2, role = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`

	now := time.Now()
	result, err := r.db.ExecContext(
		ctx, query,
		user.Name,
		user.Phone,
		user.Role,
		user.IsActive,
		now,
		user.ID,
	)

	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to update user", 500)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New(errors.ErrCodeNotFound, "User not found", 404)
	}

	user.UpdatedAt = now
	return nil
}

// UpdatePassword updates user password
func (r *repository) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, passwordHash, time.Now(), userID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to update password", 500)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New(errors.ErrCodeNotFound, "User not found", 404)
	}

	return nil
}

// UpdateLastLogin updates last login timestamp
func (r *repository) UpdateLastLogin(ctx context.Context, userID int64) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to update last login", 500)
	}

	return nil
}

// SetMFASecret sets MFA secret for user
func (r *repository) SetMFASecret(ctx context.Context, userID int64, secret string) error {
	query := `UPDATE users SET mfa_secret = $1, updated_at = $2 WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, secret, time.Now(), userID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to set MFA secret", 500)
	}

	return nil
}

// EnableMFA enables MFA for user
func (r *repository) EnableMFA(ctx context.Context, userID int64) error {
	query := `UPDATE users SET mfa_enabled = true, updated_at = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to enable MFA", 500)
	}

	return nil
}

// DisableMFA disables MFA for user
func (r *repository) DisableMFA(ctx context.Context, userID int64) error {
	query := `UPDATE users SET mfa_enabled = false, mfa_secret = NULL, updated_at = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to disable MFA", 500)
	}

	return nil
}
