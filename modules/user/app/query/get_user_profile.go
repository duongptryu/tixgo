package query

import (
	"context"

	"tixgo/modules/user/domain"
	"tixgo/shared/syserr"
)

// GetUserProfileQuery represents the query to get user profile
type GetUserProfileQuery struct {
	UserID int64 
}

// UserProfileResult represents the user profile result
type UserProfileResult struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone,omitempty"`
	UserType      string `json:"user_type"`
	Status        string `json:"status"`
	EmailVerified bool   `json:"email_verified"`
	CreatedAt     string `json:"created_at"`
	LastLogin     string `json:"last_login,omitempty"`
}

// GetUserProfileHandler handles getting user profile
type GetUserProfileHandler struct {
	userRepo domain.UserRepository
}

// NewGetUserProfileHandler creates a new get user profile handler
func NewGetUserProfileHandler(userRepo domain.UserRepository) *GetUserProfileHandler {
	return &GetUserProfileHandler{
		userRepo: userRepo,
	}
}

// Handle executes the get user profile query
func (h *GetUserProfileHandler) Handle(ctx context.Context, query GetUserProfileQuery) (*UserProfileResult, error) {
	// Get user by ID
	user, err := h.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get user")
	}

	// Convert to result
	result := &UserProfileResult{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		UserType:      string(user.UserType),
		Status:        string(user.Status),
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if user.Phone != nil {
		result.Phone = *user.Phone
	}

	if user.LastLogin != nil {
		result.LastLogin = user.LastLogin.Format("2006-01-02T15:04:05Z")
	}

	return result, nil
}
