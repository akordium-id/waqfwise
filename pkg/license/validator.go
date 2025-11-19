package license

import (
	"errors"
	"time"
)

// Validator handles license validation
type Validator struct {
	storage Storage
}

// NewValidator creates a new license validator
func NewValidator(storage Storage) *Validator {
	return &Validator{
		storage: storage,
	}
}

// Validate validates a license key
func (v *Validator) Validate(key string) (*License, error) {
	if key == "" {
		return nil, errors.New("license key is required")
	}

	// TODO: Implement license validation logic
	// - Parse license key
	// - Verify signature
	// - Check expiration
	// - Validate features
	// - Check against storage/cache

	return nil, errors.New("not implemented")
}

// ValidateFeature checks if a specific feature is enabled in the license
func (v *Validator) ValidateFeature(key string, feature string) (bool, error) {
	license, err := v.Validate(key)
	if err != nil {
		return false, err
	}

	// TODO: Check if feature is enabled in license
	_ = license
	return false, errors.New("not implemented")
}

// IsExpired checks if a license is expired
func (v *Validator) IsExpired(license *License) bool {
	return license.ExpiresAt.Before(time.Now())
}
