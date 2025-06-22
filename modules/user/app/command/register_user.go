package command

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"tixgo/modules/user/domain"

	"github.com/duongptryu/gox/syserr"
)

// RegisterUserCommand represents the command to register a new user
type RegisterUserCommand struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserType  string `json:"-" validate:"default=customer"`
}

// RegisterUserResult represents the result of user registration
type RegisterUserResult struct {
	UserID int64  `json:"user_id"`
	OTP    string `json:"otp"`
}

// RegisterUserHandler handles user registration
type RegisterUserHandler struct {
	userRepo domain.UserRepository
	otpStore domain.OTPStore
}

// NewRegisterUserHandler creates a new register user handler
func NewRegisterUserHandler(userRepo domain.UserRepository, otpStore domain.OTPStore) *RegisterUserHandler {
	return &RegisterUserHandler{
		userRepo: userRepo,
		otpStore: otpStore,
	}
}

// Handle executes the register user command
func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResult, error) {
	// Validate user type
	if !domain.IsValidUserType(cmd.UserType) {
		return nil, domain.ErrInvalidUserType
	}

	// Check if user already exists
	existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to check existing user")
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Create new user
	user, err := domain.NewUser(cmd.Email, cmd.Password, cmd.FirstName, cmd.LastName, domain.UserType(cmd.UserType))
	if err != nil {
		return nil, err
	}

	// Save user to repository
	err = h.userRepo.Create(ctx, user)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to create user")
	}

	// Generate OTP
	otp, err := generateOTP()
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to generate OTP")
	}

	// Store OTP
	err = h.otpStore.Store(ctx, cmd.Email, otp)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to store OTP")
	}

	return &RegisterUserResult{
		UserID: user.ID,
		OTP:    otp,
	}, nil
}

// generateOTP generates a 6-digit OTP
func generateOTP() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Add(n, min).Int64()), nil
}
