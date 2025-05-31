package app

import (
	"context"

	"tixgo/modules/user/app/command"
	"tixgo/modules/user/app/query"
	"tixgo/modules/user/domain"
	"tixgo/shared/auth"
)

// UserService provides the application service for user operations
type UserService struct {
	registerUserHandler   *command.RegisterUserHandler
	verifyOTPHandler      *command.VerifyOTPHandler
	loginUserHandler      *command.LoginUserHandler
	getUserProfileHandler *query.GetUserProfileHandler
}

// NewUserService creates a new user service
func NewUserService(
	userRepo domain.UserRepository,
	otpStore domain.OTPStore,
	jwtService *auth.JWTService,
) *UserService {
	return &UserService{
		registerUserHandler:   command.NewRegisterUserHandler(userRepo, otpStore),
		verifyOTPHandler:      command.NewVerifyOTPHandler(userRepo, otpStore),
		loginUserHandler:      command.NewLoginUserHandler(userRepo, jwtService),
		getUserProfileHandler: query.NewGetUserProfileHandler(userRepo),
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(ctx context.Context, cmd command.RegisterUserCommand) (*command.RegisterUserResult, error) {
	return s.registerUserHandler.Handle(ctx, cmd)
}

// VerifyOTP verifies user's OTP
func (s *UserService) VerifyOTP(ctx context.Context, cmd command.VerifyOTPCommand) (*command.VerifyOTPResult, error) {
	return s.verifyOTPHandler.Handle(ctx, cmd)
}

// LoginUser authenticates a user
func (s *UserService) LoginUser(ctx context.Context, cmd command.LoginUserCommand) (*command.LoginUserResult, error) {
	return s.loginUserHandler.Handle(ctx, cmd)
}

// GetUserProfile gets user profile
func (s *UserService) GetUserProfile(ctx context.Context, query query.GetUserProfileQuery) (*query.UserProfileResult, error) {
	return s.getUserProfileHandler.Handle(ctx, query)
}
