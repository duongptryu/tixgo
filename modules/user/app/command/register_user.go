package command

import (
	"context"

	"tixgo/modules/user/domain"

	"github.com/duongptryu/gox/messaging"
	"github.com/duongptryu/gox/syserr"
)

// RegisterUserCommand represents the command to register a new user
type RegisterUserCommand struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	UserType  string `json:"-"`
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
	eventBus      messaging.EventBus
}

// NewRegisterUserHandler creates a new register user handler
func NewRegisterUserHandler(userRepo domain.UserRepository, tempUserStore domain.TempUserStore, otpStore domain.OTPStore, eventBus messaging.EventBus) *RegisterUserHandler {
	return &RegisterUserHandler{
		userRepo:      userRepo,
		tempUserStore: tempUserStore,
		otpStore:      otpStore,
		eventBus:      eventBus,
	}
}

// Handle executes the register user command
func (h *RegisterUserHandler) Handle(ctx context.Context, cmd *RegisterUserCommand) (*RegisterUserResult, error) {
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
	user, err := domain.NewUserCustomer(cmd.Email, cmd.Password, cmd.FirstName, cmd.LastName)
	if err != nil {
		return nil, err
	}

	// Store user temporarily (not in database yet)
	err = h.tempUserStore.Store(ctx, cmd.Email, user)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to store user temporarily")
	}

	// Publish event to send OTP to user
	err = h.eventBus.PublishEvent(ctx, domain.NewEventUserRegistered(user.Email))
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to publish event user registered")
	}

	return &RegisterUserResult{
		Email: user.Email,
	}, nil
}
