package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/akordium-id/waqfwise/internal/services/auth/dto"
	"github.com/akordium-id/waqfwise/internal/services/auth/service"
	"github.com/akordium-id/waqfwise/internal/shared/errors"
	"github.com/akordium-id/waqfwise/internal/shared/response"
	"github.com/akordium-id/waqfwise/internal/shared/validator"
	"github.com/gorilla/mux"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	service service.Service
}

// New creates a new authentication handler
func New(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.New(errors.ErrCodeBadRequest, "Invalid request body", 400))
		return
	}

	// Validate request
	v := validator.New()
	v.Required("email", req.Email)
	v.Email("email", req.Email)
	v.Required("password", req.Password)
	v.Password("password", req.Password)
	v.Required("name", req.Name)
	v.MinLength("name", req.Name, 2)

	if req.Phone != "" {
		v.Phone("phone", req.Phone)
	}

	if !v.IsValid() {
		response.Error(w, v.Error())
		return
	}

	// Register user
	authResp, err := h.service.Register(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, authResp)
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.New(errors.ErrCodeBadRequest, "Invalid request body", 400))
		return
	}

	// Validate request
	v := validator.New()
	v.Required("email", req.Email)
	v.Email("email", req.Email)
	v.Required("password", req.Password)

	if !v.IsValid() {
		response.Error(w, v.Error())
		return
	}

	// Login
	authResp, err := h.service.Login(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, authResp)
}

// RefreshToken handles token refresh
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.New(errors.ErrCodeBadRequest, "Invalid request body", 400))
		return
	}

	if req.RefreshToken == "" {
		response.Error(w, errors.New(errors.ErrCodeValidation, "Refresh token is required", 400))
		return
	}

	authResp, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, authResp)
}

// GetProfile handles get user profile
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	profile, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, profile)
}

// ChangePassword handles password change
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.New(errors.ErrCodeBadRequest, "Invalid request body", 400))
		return
	}

	// Validate request
	v := validator.New()
	v.Required("old_password", req.OldPassword)
	v.Required("new_password", req.NewPassword)
	v.Password("new_password", req.NewPassword)

	if !v.IsValid() {
		response.Error(w, v.Error())
		return
	}

	if err := h.service.ChangePassword(r.Context(), userID, &req); err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, map[string]string{"message": "Password changed successfully"})
}

// SetupMFA handles MFA setup
func (h *Handler) SetupMFA(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	mfaSetup, err := h.service.SetupMFA(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, mfaSetup)
}

// EnableMFA handles MFA enable
func (h *Handler) EnableMFA(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	var req dto.EnableMFARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.New(errors.ErrCodeBadRequest, "Invalid request body", 400))
		return
	}

	if req.Code == "" {
		response.Error(w, errors.New(errors.ErrCodeValidation, "MFA code is required", 400))
		return
	}

	if err := h.service.EnableMFA(r.Context(), userID, req.Code); err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, map[string]string{"message": "MFA enabled successfully"})
}

// DisableMFA handles MFA disable
func (h *Handler) DisableMFA(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	if err := h.service.DisableMFA(r.Context(), userID); err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, map[string]string{"message": "MFA disabled successfully"})
}

// ValidateToken handles token validation
func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.Error(w, errors.ErrInvalidToken)
		return
	}

	claims, err := h.service.ValidateToken(parts[1])
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, map[string]interface{}{
		"valid":     true,
		"user_id":   claims.UserID,
		"email":     claims.Email,
		"role":      claims.Role,
		"tenant_id": claims.TenantID,
	})
}

// RegisterRoutes registers HTTP routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	// Public routes
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/refresh", h.RefreshToken).Methods("POST")
	r.HandleFunc("/validate", h.ValidateToken).Methods("POST")

	// Protected routes (require auth middleware)
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(h.authMiddleware)
	protected.HandleFunc("/profile", h.GetProfile).Methods("GET")
	protected.HandleFunc("/change-password", h.ChangePassword).Methods("POST")
	protected.HandleFunc("/mfa/setup", h.SetupMFA).Methods("POST")
	protected.HandleFunc("/mfa/enable", h.EnableMFA).Methods("POST")
	protected.HandleFunc("/mfa/disable", h.DisableMFA).Methods("POST")
}

// authMiddleware authenticates requests
func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, errors.ErrUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, errors.ErrInvalidToken)
			return
		}

		claims, err := h.service.ValidateToken(parts[1])
		if err != nil {
			response.Error(w, err)
			return
		}

		// Add user ID to context
		ctx := r.Context()
		ctx = setUserIDToContext(ctx, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Context helpers
type contextKey string

const userIDKey contextKey = "user_id"

func setUserIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func getUserIDFromContext(r *http.Request) int64 {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}
