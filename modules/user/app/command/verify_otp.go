package command

import (
	"context"

	"tixgo/modules/user/domain"

	"github.com/duongptryu/gox/syserr"
)

// VerifyOTPCommand represents the command to verify OTP
type VerifyOTPCommand struct {
	Email string `json:"email"`	
	OTP   string `json:"otp"`
}

// VerifyOTPResult represents the result of OTP verification
type VerifyOTPResult struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

// VerifyOTPHandler handles OTP verification
type VerifyOTPHandler struct {
	userRepo domain.UserRepository
	otpStore domain.OTPStore
}

// NewVerifyOTPHandler creates a new verify OTP handler
func NewVerifyOTPHandler(userRepo domain.UserRepository, otpStore domain.OTPStore) *VerifyOTPHandler {
	return &VerifyOTPHandler{
		userRepo: userRepo,
		otpStore: otpStore,
	}
}

// Handle executes the verify OTP command
func (h *VerifyOTPHandler) Handle(ctx context.Context, cmd VerifyOTPCommand) (*VerifyOTPResult, error) {
	// Verify OTP
	err := h.otpStore.Verify(ctx, cmd.Email, cmd.OTP)
	if err != nil {
		return nil, domain.ErrInvalidOTP
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get user")
	}

	// Mark email as verified
	user.VerifyEmail()

	// Update user in repository
	err = h.userRepo.Update(ctx, user)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to update user")
	}

	return &VerifyOTPResult{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}
