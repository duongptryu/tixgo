package adapters

import (
	"context"
	"testing"
	"time"

	"tixgo/modules/user/domain"
)

func TestInMemoryOTPStore_Store(t *testing.T) {
	store := NewInMemoryOTPStore()
	defer store.Close()

	ctx := context.Background()
	email := "test@example.com"
	otp := "123456"

	err := store.Store(ctx, email, otp)
	if err != nil {
		t.Errorf("Store() unexpected error = %v", err)
	}

	// Verify OTP was stored
	store.mutex.RLock()
	entry, exists := store.store[email]
	store.mutex.RUnlock()

	if !exists {
		t.Errorf("Store() OTP was not stored")
	}

	if entry.OTP != otp {
		t.Errorf("Store() stored OTP = %v, want %v", entry.OTP, otp)
	}

	if time.Until(entry.ExpiresAt) > 5*time.Minute {
		t.Errorf("Store() expiration time is too far in the future")
	}
}

func TestInMemoryOTPStore_Verify(t *testing.T) {
	store := NewInMemoryOTPStore()
	defer store.Close()

	ctx := context.Background()
	email := "test@example.com"
	otp := "123456"

	tests := []struct {
		name    string
		setup   func()
		email   string
		otp     string
		wantErr error
	}{
		{
			name: "valid OTP",
			setup: func() {
				store.Store(ctx, email, otp)
			},
			email:   email,
			otp:     otp,
			wantErr: nil,
		},
		{
			name:    "non-existent OTP",
			setup:   func() {},
			email:   "nonexistent@example.com",
			otp:     "123456",
			wantErr: domain.ErrInvalidOTP,
		},
		{
			name: "wrong OTP",
			setup: func() {
				store.Store(ctx, email, otp)
			},
			email:   email,
			otp:     "wrong",
			wantErr: domain.ErrInvalidOTP,
		},
		{
			name: "expired OTP",
			setup: func() {
				// Manually create expired entry
				store.mutex.Lock()
				store.store[email] = &OTPEntry{
					OTP:       otp,
					ExpiresAt: time.Now().Add(-1 * time.Minute),
				}
				store.mutex.Unlock()
			},
			email:   email,
			otp:     otp,
			wantErr: domain.ErrOTPExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear store
			store.mutex.Lock()
			store.store = make(map[string]*OTPEntry)
			store.mutex.Unlock()

			// Setup test case
			tt.setup()

			err := store.Verify(ctx, tt.email, tt.otp)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Verify() expected error %v but got none", tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("Verify() error = %v, want %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Verify() unexpected error = %v", err)
				}

				// Verify OTP was removed after successful verification
				store.mutex.RLock()
				_, exists := store.store[tt.email]
				store.mutex.RUnlock()

				if exists {
					t.Errorf("Verify() OTP should be removed after successful verification")
				}
			}
		})
	}
}

func TestInMemoryOTPStore_Delete(t *testing.T) {
	store := NewInMemoryOTPStore()
	defer store.Close()

	ctx := context.Background()
	email := "test@example.com"
	otp := "123456"

	// Store an OTP
	err := store.Store(ctx, email, otp)
	if err != nil {
		t.Fatalf("Failed to store OTP: %v", err)
	}

	// Verify it exists
	store.mutex.RLock()
	_, exists := store.store[email]
	store.mutex.RUnlock()

	if !exists {
		t.Fatalf("OTP should exist before deletion")
	}

	// Delete the OTP
	err = store.Delete(ctx, email)
	if err != nil {
		t.Errorf("Delete() unexpected error = %v", err)
	}

	// Verify it was deleted
	store.mutex.RLock()
	_, exists = store.store[email]
	store.mutex.RUnlock()

	if exists {
		t.Errorf("Delete() OTP should be removed after deletion")
	}
}

func TestInMemoryOTPStore_CleanupExpired(t *testing.T) {
	store := NewInMemoryOTPStore()
	defer store.Close()

	ctx := context.Background()

	// Add valid OTP
	validEmail := "valid@example.com"
	validOTP := "123456"
	store.Store(ctx, validEmail, validOTP)

	// Add expired OTP manually
	expiredEmail := "expired@example.com"
	store.mutex.Lock()
	store.store[expiredEmail] = &OTPEntry{
		OTP:       "654321",
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	store.mutex.Unlock()

	// Run cleanup
	store.cleanupExpired()

	// Verify valid OTP still exists
	store.mutex.RLock()
	_, validExists := store.store[validEmail]
	_, expiredExists := store.store[expiredEmail]
	store.mutex.RUnlock()

	if !validExists {
		t.Errorf("cleanupExpired() should not remove valid OTP")
	}

	if expiredExists {
		t.Errorf("cleanupExpired() should remove expired OTP")
	}
}
