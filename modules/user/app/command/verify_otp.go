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
	userRepo      domain.UserRepository
	tempUserStore domain.TempUserStore
	otpStore      domain.OTPStore
}

// NewVerifyOTPHandler creates a new verify OTP handler
func NewVerifyOTPHandler(userRepo domain.UserRepository, tempUserStore domain.TempUserStore, otpStore domain.OTPStore) *VerifyOTPHandler {
	return &VerifyOTPHandler{
		userRepo:      userRepo,
		tempUserStore: tempUserStore,
		otpStore:      otpStore,
	}
}

// Handle executes the verify OTP command
func (h *VerifyOTPHandler) Handle(ctx context.Context, cmd *VerifyOTPCommand) (*VerifyOTPResult, error) {
	// Verify OTP
	err := h.otpStore.Verify(ctx, cmd.Email, cmd.OTP)
	if err != nil {
		return nil, domain.ErrInvalidOTP
	}

	// Get user from temp store
	user, err := h.tempUserStore.Get(ctx, cmd.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get temp user")
	}

	// Mark email as verified
	user.VerifyEmail()

	// Save user to database (move from temp to permanent storage)
	err = h.userRepo.Create(ctx, user)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to create user")
	}

	// Clean up temp store
	err = h.tempUserStore.Delete(ctx, cmd.Email)
	if err != nil {
		// Log error but don't fail the operation since user is already created
		// This is just cleanup
	}

	return &VerifyOTPResult{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}
