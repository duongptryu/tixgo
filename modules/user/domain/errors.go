package domain

import "tixgo/shared/syserr"

// Domain-specific error codes for client handling
const (
	// User not found errors
	UserNotFoundCode syserr.Code = "user_not_found"

	// User registration errors
	UserAlreadyExistsCode syserr.Code = "user_already_exists"
	InvalidUserTypeCode   syserr.Code = "invalid_user_type"

	// Authentication errors
	InvalidCredentialsCode syserr.Code = "invalid_credentials"

	// Authorization/Access errors
	EmailNotVerifiedCode syserr.Code = "email_not_verified"
	UserInactiveCode     syserr.Code = "user_inactive"
	UserSuspendedCode    syserr.Code = "user_suspended"

	// OTP errors
	InvalidOTPCode  syserr.Code = "invalid_otp"
	OTPExpiredCode  syserr.Code = "otp_expired"
	OTPNotFoundCode syserr.Code = "otp_not_found"
)

// Domain-specific errors with specific codes
var (
	// User not found errors
	ErrUserNotFound = syserr.New(UserNotFoundCode, "user not found")

	// User registration errors
	ErrUserAlreadyExists = syserr.New(UserAlreadyExistsCode, "user with this email already exists")
	ErrInvalidUserType   = syserr.New(InvalidUserTypeCode, "invalid user type, must be: customer, organizer, or admin")

	// Authentication errors
	ErrInvalidCredentials = syserr.New(InvalidCredentialsCode, "invalid email or password")

	// Authorization/Access errors
	ErrEmailNotVerified = syserr.New(EmailNotVerifiedCode, "email address not verified, please check your email for verification code")
	ErrUserInactive     = syserr.New(UserInactiveCode, "user account is inactive, please contact support")
	ErrUserSuspended    = syserr.New(UserSuspendedCode, "user account is suspended, please contact support")

	// OTP errors
	ErrInvalidOTP  = syserr.New(InvalidOTPCode, "invalid verification code")
	ErrOTPExpired  = syserr.New(OTPExpiredCode, "verification code has expired, please request a new one")
	ErrOTPNotFound = syserr.New(OTPNotFoundCode, "no verification code found for this email")
)
