package command

import (
	"context"
	"strconv"

	"tixgo/modules/user/domain"
	"tixgo/shared/auth"
	"tixgo/shared/syserr"
)

// LoginUserCommand represents the command to login a user
type LoginUserCommand struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUserResult represents the result of user login
type LoginUserResult struct {
	UserID       int64  `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// LoginUserHandler handles user login
type LoginUserHandler struct {
	userRepo   domain.UserRepository
	jwtService *auth.JWTService
}

// NewLoginUserHandler creates a new login user handler
func NewLoginUserHandler(userRepo domain.UserRepository, jwtService *auth.JWTService) *LoginUserHandler {
	return &LoginUserHandler{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Handle executes the login user command
func (h *LoginUserHandler) Handle(ctx context.Context, cmd LoginUserCommand) (*LoginUserResult, error) {
	// Get user by email
	user, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get user")
	}

	// Check password
	err = user.CheckPassword(cmd.Password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check if user can login
	err = user.CanLogin()
	if err != nil {
		return nil, err
	}

	// Update last login
	user.UpdateLastLogin()
	err = h.userRepo.Update(ctx, user)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to update last login")
	}

	// Generate JWT tokens
	accessToken, refreshToken, expiresIn, err := h.jwtService.GenerateTokenPair(ctx, strconv.FormatInt(user.ID, 10), string(user.UserType))
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to generate tokens")
	}

	return &LoginUserResult{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
