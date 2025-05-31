package domain

import "time"

// DomainEvent represents a domain event
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

// UserRegistered event is published when a user registers
type UserRegistered struct {
	UserID     int64
	Email      string
	UserType   UserType
	occurredAt time.Time
}

func NewUserRegistered(userID int64, email string, userType UserType) *UserRegistered {
	return &UserRegistered{
		UserID:     userID,
		Email:      email,
		UserType:   userType,
		occurredAt: time.Now(),
	}
}

func (e *UserRegistered) EventType() string {
	return "user.registered"
}

func (e *UserRegistered) OccurredAt() time.Time {
	return e.occurredAt
}

// UserEmailVerified event is published when a user verifies their email
type UserEmailVerified struct {
	UserID     int64
	Email      string
	occurredAt time.Time
}

func NewUserEmailVerified(userID int64, email string) *UserEmailVerified {
	return &UserEmailVerified{
		UserID:     userID,
		Email:      email,
		occurredAt: time.Now(),
	}
}

func (e *UserEmailVerified) EventType() string {
	return "user.email_verified"
}

func (e *UserEmailVerified) OccurredAt() time.Time {
	return e.occurredAt
}

// UserLoggedIn event is published when a user logs in
type UserLoggedIn struct {
	UserID     int64
	Email      string
	occurredAt time.Time
}

func NewUserLoggedIn(userID int64, email string) *UserLoggedIn {
	return &UserLoggedIn{
		UserID:     userID,
		Email:      email,
		occurredAt: time.Now(),
	}
}

func (e *UserLoggedIn) EventType() string {
	return "user.logged_in"
}

func (e *UserLoggedIn) OccurredAt() time.Time {
	return e.occurredAt
}
