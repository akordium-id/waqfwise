package validator

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/akordium-id/waqfwise/internal/shared/errors"
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validator provides validation functions
type Validator struct {
	errors []ValidationError
}

// New creates a new validator
func New() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// IsValid returns true if there are no validation errors
func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

// Errors returns all validation errors
func (v *Validator) Errors() []ValidationError {
	return v.errors
}

// Error returns validation errors as AppError
func (v *Validator) Error() *errors.AppError {
	if v.IsValid() {
		return nil
	}

	messages := make([]string, len(v.errors))
	for i, err := range v.errors {
		messages[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
	}

	return errors.New(
		errors.ErrCodeValidation,
		strings.Join(messages, "; "),
		400,
	)
}

// AddError adds a validation error
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// Required validates that a field is not empty
func (v *Validator) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "is required")
	}
}

// Email validates email format
func (v *Validator) Email(field, email string) {
	if email == "" {
		return
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		v.AddError(field, "must be a valid email address")
	}
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field, value string, min int) {
	if len(value) < min {
		v.AddError(field, fmt.Sprintf("must be at least %d characters", min))
	}
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field, value string, max int) {
	if len(value) > max {
		v.AddError(field, fmt.Sprintf("must be at most %d characters", max))
	}
}

// Min validates minimum numeric value
func (v *Validator) Min(field string, value, min int64) {
	if value < min {
		v.AddError(field, fmt.Sprintf("must be at least %d", min))
	}
}

// Max validates maximum numeric value
func (v *Validator) Max(field string, value, max int64) {
	if value > max {
		v.AddError(field, fmt.Sprintf("must be at most %d", max))
	}
}

// In validates that value is in allowed list
func (v *Validator) In(field, value string, allowed []string) {
	if value == "" {
		return
	}

	for _, a := range allowed {
		if value == a {
			return
		}
	}

	v.AddError(field, fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")))
}

// Matches validates that value matches regex pattern
func (v *Validator) Matches(field, value, pattern string) {
	if value == "" {
		return
	}

	matched, err := regexp.MatchString(pattern, value)
	if err != nil || !matched {
		v.AddError(field, "has invalid format")
	}
}

// Phone validates Indonesian phone number format
func (v *Validator) Phone(field, phone string) {
	if phone == "" {
		return
	}

	// Indonesian phone number pattern
	pattern := `^(\+62|62|0)[0-9]{9,12}$`
	v.Matches(field, phone, pattern)
}

// Password validates password strength
func (v *Validator) Password(field, password string) {
	if password == "" {
		return
	}

	if len(password) < 8 {
		v.AddError(field, "must be at least 8 characters")
		return
	}

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	)

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		v.AddError(field, "must contain uppercase, lowercase, number, and special character")
	}
}
