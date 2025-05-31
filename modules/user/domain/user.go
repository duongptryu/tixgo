package domain

import (
	"time"

	"tixgo/shared/syserr"

	"golang.org/x/crypto/bcrypt"
)

// UserType represents the type of user
type UserType string

const (
	UserTypeCustomer  UserType = "customer"
	UserTypeOrganizer UserType = "organizer"
	UserTypeAdmin     UserType = "admin"
)

// UserStatus represents the status of user
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// User represents the user aggregate root
type User struct {
	ID            int64
	Email         string
	PasswordHash  string
	FirstName     string
	LastName      string
	Phone         *string
	DateOfBirth   *time.Time
	UserType      UserType
	Status        UserStatus
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastLogin     *time.Time
}

// NewUser creates a new user with hashed password
func NewUser(email, password, firstName, lastName string, userType UserType) (*User, error) {
	if email == "" {
		return nil, syserr.New(syserr.InvalidArgumentCode, "email is required")
	}
	if password == "" {
		return nil, syserr.New(syserr.InvalidArgumentCode, "password is required")
	}
	if firstName == "" {
		return nil, syserr.New(syserr.InvalidArgumentCode, "first name is required")
	}
	if lastName == "" {
		return nil, syserr.New(syserr.InvalidArgumentCode, "last name is required")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to hash password")
	}

	now := time.Now()
	return &User{
		Email:         email,
		PasswordHash:  hashedPassword,
		FirstName:     firstName,
		LastName:      lastName,
		UserType:      userType,
		Status:        UserStatusActive,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// CheckPassword verifies if the provided password matches the user's password
func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return syserr.New(syserr.UnauthorizedCode, "invalid password")
	}
	return nil
}

// VerifyEmail marks the user's email as verified
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now()
}

// UpdateLastLogin updates the user's last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

// CanLogin checks if the user can login
func (u *User) CanLogin() error {
	if u.Status != UserStatusActive {
		return syserr.New(syserr.ForbiddenCode, "user account is not active")
	}
	if !u.EmailVerified {
		return syserr.New(syserr.ForbiddenCode, "email not verified")
	}
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// hashPassword hashes the password using bcrypt
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// IsValidUserType checks if the user type is valid
func IsValidUserType(userType string) bool {
	switch UserType(userType) {
	case UserTypeCustomer, UserTypeOrganizer, UserTypeAdmin:
		return true
	default:
		return false
	}
}
