package adapters

import (
	"context"
	"sync"
	"time"

	"tixgo/modules/user/domain"
)

// OTPEntry represents an OTP entry with expiration
type OTPEntry struct {
	OTP       string
	ExpiresAt time.Time
}

// InMemoryOTPStore implements the OTPStore interface using in-memory storage
type InMemoryOTPStore struct {
	store   map[string]*OTPEntry
	mutex   sync.RWMutex
	cleanup chan struct{}
}

// NewInMemoryOTPStore creates a new in-memory OTP store
func NewInMemoryOTPStore() *InMemoryOTPStore {
	store := &InMemoryOTPStore{
		store:   make(map[string]*OTPEntry),
		cleanup: make(chan struct{}),
	}

	// Start cleanup goroutine
	go store.startCleanup()

	return store
}

// Store stores an OTP for a user email with 5-minute expiration
func (s *InMemoryOTPStore) Store(ctx context.Context, email, otp string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[email] = &OTPEntry{
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return nil
}

// Verify verifies an OTP for a user email and removes it if valid
func (s *InMemoryOTPStore) Verify(ctx context.Context, email, otp string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	entry, exists := s.store[email]
	if !exists {
		return domain.ErrInvalidOTP
	}

	// Check if OTP has expired
	if time.Now().After(entry.ExpiresAt) {
		delete(s.store, email)
		return domain.ErrOTPExpired
	}

	// Check if OTP matches
	if entry.OTP != otp {
		return domain.ErrInvalidOTP
	}

	// Remove OTP after successful verification
	delete(s.store, email)

	return nil
}

// Delete removes an OTP for a user email
func (s *InMemoryOTPStore) Delete(ctx context.Context, email string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, email)
	return nil
}

// startCleanup starts a goroutine to clean up expired OTPs
func (s *InMemoryOTPStore) startCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
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

// cleanupExpired removes expired OTPs from the store
func (s *InMemoryOTPStore) cleanupExpired() {
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
func (s *InMemoryOTPStore) Close() {
	close(s.cleanup)
}
