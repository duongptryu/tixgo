package command

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"tixgo/modules/user/domain"

	"github.com/duongptryu/gox/eventbus"
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
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

// RegisterUserHandler handles user registration
type RegisterUserHandler struct {
	userRepo      domain.UserRepository
	tempUserStore domain.TempUserStore
	otpStore      domain.OTPStore
	commandBus    eventbus.BusCommand
}

// NewRegisterUserHandler creates a new register user handler
func NewRegisterUserHandler(userRepo domain.UserRepository, tempUserStore domain.TempUserStore, otpStore domain.OTPStore, commandBus eventbus.BusCommand) *RegisterUserHandler {
	return &RegisterUserHandler{
		userRepo:      userRepo,
		tempUserStore: tempUserStore,
		otpStore:      otpStore,
		commandBus:    commandBus,
	}
}

// Handle executes the register user command
func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResult, error) {
	// Validate user type
	if !domain.IsValidUserType(cmd.UserType) {
		return nil, domain.ErrInvalidUserType
	}

	// Check if user already exists in database
	existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to check existing user")
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Check if user already exists in temp store (pending verification)
	tempUser, err := h.tempUserStore.Get(ctx, cmd.Email)
	if err == nil && tempUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Create new user
	user, err := domain.NewUser(cmd.Email, cmd.Password, cmd.FirstName, cmd.LastName, domain.UserType(cmd.UserType))
	if err != nil {
		return nil, err
	}

	// Store user temporarily (not in database yet)
	err = h.tempUserStore.Store(ctx, cmd.Email, user)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to store user temporarily")
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

	// Publish event to send OTP to user
	err = h.commandBus.PublishCommand(ctx, domain.NewCommandSendUserMailOTP(user.Email, otp))
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to publish event send user mail otp")
	}

	return &RegisterUserResult{
		Email: user.Email,
		OTP:   otp,
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
