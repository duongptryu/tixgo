package adapters

import (
	"context"
	"sync"
	"time"

	"tixgo/modules/user/domain"
)

// TempUserEntry represents a temporary user entry with expiration
type TempUserEntry struct {
	User      *domain.User
	ExpiresAt time.Time
}

// InMemoryTempUserStore implements the TempUserStore interface using in-memory storage
type InMemoryTempUserStore struct {
	store   map[string]*TempUserEntry
	mutex   sync.RWMutex
	cleanup chan struct{}
}

// NewInMemoryTempUserStore creates a new in-memory temporary user store
func NewInMemoryTempUserStore() *InMemoryTempUserStore {
	store := &InMemoryTempUserStore{
		store:   make(map[string]*TempUserEntry),
		cleanup: make(chan struct{}),
	}

	// Start cleanup goroutine
	go store.startCleanup()

	return store
}

// Store stores a user temporarily with 10-minute expiration
func (s *InMemoryTempUserStore) Store(ctx context.Context, email string, user *domain.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[email] = &TempUserEntry{
		User:      user,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	return nil
}

// Get retrieves a temporary user by email
func (s *InMemoryTempUserStore) Get(ctx context.Context, email string) (*domain.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entry, exists := s.store[email]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	// Check if entry has expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, domain.ErrUserNotFound
	}

	return entry.User, nil
}

// Delete removes a temporary user by email
func (s *InMemoryTempUserStore) Delete(ctx context.Context, email string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, email)
	return nil
}

// startCleanup starts a goroutine to clean up expired temporary users
func (s *InMemoryTempUserStore) startCleanup() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupExpired()
		case <-s.cleanup:
			return
		}
	}
}

// cleanupExpired removes expired temporary users from the store
func (s *InMemoryTempUserStore) cleanupExpired() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for email, entry := range s.store {
		if now.After(entry.ExpiresAt) {
			delete(s.store, email)
		}
	}
}

// Close stops the cleanup goroutine
func (s *InMemoryTempUserStore) Close() {
	close(s.cleanup)
}
