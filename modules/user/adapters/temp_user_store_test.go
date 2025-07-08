package adapters

import (
	"context"
	"testing"
	"time"

	"tixgo/modules/user/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryTempUserStore_Store(t *testing.T) {
	store := NewInMemoryTempUserStore()
	defer store.Close()

	ctx := context.Background()
	email := "test@example.com"

	user, err := domain.NewUserCustomer(email, "password123", "John", "Doe")
	require.NoError(t, err)

	err = store.Store(ctx, email, user)
	assert.NoError(t, err)

	// Verify user was stored
	retrievedUser, err := store.Get(ctx, email)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.FirstName, retrievedUser.FirstName)
	assert.Equal(t, user.LastName, retrievedUser.LastName)
}

func TestInMemoryTempUserStore_Get_NotFound(t *testing.T) {
	store := NewInMemoryTempUserStore()
	defer store.Close()

	ctx := context.Background()
	email := "nonexistent@example.com"

	_, err := store.Get(ctx, email)
	assert.Equal(t, domain.ErrUserNotFound, err)
}

func TestInMemoryTempUserStore_Delete(t *testing.T) {
	store := NewInMemoryTempUserStore()
	defer store.Close()

	ctx := context.Background()
	email := "test@example.com"

	user, err := domain.NewUserCustomer(email, "password123", "John", "Doe")
	require.NoError(t, err)

	// Store user
	err = store.Store(ctx, email, user)
	require.NoError(t, err)

	// Verify user exists
	_, err = store.Get(ctx, email)
	assert.NoError(t, err)

	// Delete user
	err = store.Delete(ctx, email)
	assert.NoError(t, err)

	// Verify user is deleted
	_, err = store.Get(ctx, email)
	assert.Equal(t, domain.ErrUserNotFound, err)
}

func TestInMemoryTempUserStore_Expiration(t *testing.T) {
	store := NewInMemoryTempUserStore()
	defer store.Close()

	ctx := context.Background()
	email := "test@example.com"

	user, err := domain.NewUserCustomer(email, "password123", "John", "Doe")
	require.NoError(t, err)

	// Manually set a very short expiration for testing
	store.mutex.Lock()
	store.store[email] = &TempUserEntry{
		User:      user,
		ExpiresAt: time.Now().Add(1 * time.Millisecond),
	}
	store.mutex.Unlock()

	// Wait for expiration
	time.Sleep(2 * time.Millisecond)

	// Verify user is expired
	_, err = store.Get(ctx, email)
	assert.Equal(t, domain.ErrUserNotFound, err)
}

func TestInMemoryTempUserStore_CleanupExpired(t *testing.T) {
	store := NewInMemoryTempUserStore()
	defer store.Close()

	ctx := context.Background()
	email1 := "test1@example.com"
	email2 := "test2@example.com"

	user1, err := domain.NewUserCustomer(email1, "password123", "John", "Doe")
	require.NoError(t, err)

	user2, err := domain.NewUserCustomer(email2, "password123", "Jane", "Smith")
	require.NoError(t, err)

	// Store both users with different expirations
	store.mutex.Lock()
	store.store[email1] = &TempUserEntry{
		User:      user1,
		ExpiresAt: time.Now().Add(-1 * time.Minute), // Already expired
	}
	store.store[email2] = &TempUserEntry{
		User:      user2,
		ExpiresAt: time.Now().Add(10 * time.Minute), // Not expired
	}
	store.mutex.Unlock()

	// Run cleanup
	store.cleanupExpired()

	// Check that expired user is removed
	_, err = store.Get(ctx, email1)
	assert.Equal(t, domain.ErrUserNotFound, err)

	// Check that non-expired user still exists
	_, err = store.Get(ctx, email2)
	assert.NoError(t, err)
}
