package license

import "errors"

// Storage defines the interface for license storage/caching
type Storage interface {
	// Get retrieves a license from storage
	Get(key string) (*License, error)

	// Set stores a license in storage
	Set(key string, license *License) error

	// Delete removes a license from storage
	Delete(key string) error

	// Exists checks if a license exists in storage
	Exists(key string) (bool, error)
}

// MemoryStorage implements in-memory storage for licenses
type MemoryStorage struct {
	licenses map[string]*License
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		licenses: make(map[string]*License),
	}
}

// Get retrieves a license from memory
func (s *MemoryStorage) Get(key string) (*License, error) {
	license, ok := s.licenses[key]
	if !ok {
		return nil, errors.New("license not found")
	}
	return license, nil
}

// Set stores a license in memory
func (s *MemoryStorage) Set(key string, license *License) error {
	s.licenses[key] = license
	return nil
}

// Delete removes a license from memory
func (s *MemoryStorage) Delete(key string) error {
	delete(s.licenses, key)
	return nil
}

// Exists checks if a license exists in memory
func (s *MemoryStorage) Exists(key string) (bool, error) {
	_, ok := s.licenses[key]
	return ok, nil
}

// TODO: Implement RedisStorage for production use
