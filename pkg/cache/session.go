package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Session represents a user session
type Session struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt time.Time              `json:"expires_at"`
}

// SessionStore manages user sessions in Redis
type SessionStore struct {
	redis  *Redis
	prefix string
	ttl    time.Duration
}

// NewSessionStore creates a new session store
func NewSessionStore(redis *Redis, ttl time.Duration) *SessionStore {
	return &SessionStore{
		redis:  redis,
		prefix: "session:",
		ttl:    ttl,
	}
}

// Create creates a new session
func (s *SessionStore) Create(ctx context.Context, userID string, data map[string]interface{}) (*Session, error) {
	sessionID := uuid.New().String()
	now := time.Now()

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		Data:      data,
		CreatedAt: now,
		ExpiresAt: now.Add(s.ttl),
	}

	key := s.getKey(sessionID)
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := s.redis.Set(ctx, key, sessionJSON, s.ttl); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// Get retrieves a session by ID
func (s *SessionStore) Get(ctx context.Context, sessionID string) (*Session, error) {
	key := s.getKey(sessionID)
	sessionJSON, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	var session Session
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		_ = s.Delete(ctx, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return &session, nil
}

// Update updates a session
func (s *SessionStore) Update(ctx context.Context, sessionID string, data map[string]interface{}) error {
	session, err := s.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	session.Data = data
	key := s.getKey(sessionID)
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Refresh TTL
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		ttl = s.ttl
		session.ExpiresAt = time.Now().Add(ttl)
	}

	if err := s.redis.Set(ctx, key, sessionJSON, ttl); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// Delete deletes a session
func (s *SessionStore) Delete(ctx context.Context, sessionID string) error {
	key := s.getKey(sessionID)
	return s.redis.Delete(ctx, key)
}

// Refresh refreshes a session's TTL
func (s *SessionStore) Refresh(ctx context.Context, sessionID string) error {
	session, err := s.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(s.ttl)
	key := s.getKey(sessionID)
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := s.redis.Set(ctx, key, sessionJSON, s.ttl); err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	return nil
}

func (s *SessionStore) getKey(sessionID string) string {
	return s.prefix + sessionID
}
