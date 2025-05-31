package domain

import "context"

// UserRepository defines the interface for user persistence
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id int64) error
}

// OTPStore defines the interface for OTP storage and verification
type OTPStore interface {
	// Store stores an OTP for a user email with expiration
	Store(ctx context.Context, email, otp string) error

	// Verify verifies an OTP for a user email and removes it if valid
	Verify(ctx context.Context, email, otp string) error

	// Delete removes an OTP for a user email
	Delete(ctx context.Context, email string) error
}
