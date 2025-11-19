package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

// Error implements error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap implements error unwrapping
func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	ErrCodeInternal        = "INTERNAL_ERROR"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeBadRequest      = "BAD_REQUEST"
	ErrCodeConflict        = "CONFLICT"
	ErrCodeValidation      = "VALIDATION_ERROR"
	ErrCodeInvalidToken    = "INVALID_TOKEN"
	ErrCodeExpiredToken    = "EXPIRED_TOKEN"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeDuplicateEntry  = "DUPLICATE_ENTRY"
	ErrCodePaymentFailed   = "PAYMENT_FAILED"
	ErrCodeInsufficientFunds = "INSUFFICIENT_FUNDS"
	ErrCodeFraudDetected   = "FRAUD_DETECTED"
)

// New creates a new AppError
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Wrap wraps an existing error with AppError
func Wrap(err error, code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// Common errors
var (
	ErrInternal = New(ErrCodeInternal, "Internal server error", http.StatusInternalServerError)
	ErrNotFound = New(ErrCodeNotFound, "Resource not found", http.StatusNotFound)
	ErrUnauthorized = New(ErrCodeUnauthorized, "Unauthorized access", http.StatusUnauthorized)
	ErrForbidden = New(ErrCodeForbidden, "Forbidden access", http.StatusForbidden)
	ErrBadRequest = New(ErrCodeBadRequest, "Bad request", http.StatusBadRequest)
	ErrConflict = New(ErrCodeConflict, "Resource conflict", http.StatusConflict)
	ErrInvalidToken = New(ErrCodeInvalidToken, "Invalid or malformed token", http.StatusUnauthorized)
	ErrExpiredToken = New(ErrCodeExpiredToken, "Token has expired", http.StatusUnauthorized)
	ErrInvalidCredentials = New(ErrCodeInvalidCredentials, "Invalid email or password", http.StatusUnauthorized)
	ErrDuplicateEntry = New(ErrCodeDuplicateEntry, "Duplicate entry found", http.StatusConflict)
	ErrPaymentFailed = New(ErrCodePaymentFailed, "Payment processing failed", http.StatusPaymentRequired)
	ErrFraudDetected = New(ErrCodeFraudDetected, "Transaction flagged as suspicious", http.StatusForbidden)
)

// IsNotFound checks if error is not found error
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeNotFound
	}
	return false
}

// IsUnauthorized checks if error is unauthorized error
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeUnauthorized || appErr.Code == ErrCodeInvalidToken || appErr.Code == ErrCodeExpiredToken
	}
	return false
}

// IsValidationError checks if error is validation error
func IsValidationError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeValidation
	}
	return false
}

// GetHTTPStatus extracts HTTP status from error
func GetHTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Status
	}
	return http.StatusInternalServerError
}

// GetErrorCode extracts error code from error
func GetErrorCode(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return ErrCodeInternal
}
