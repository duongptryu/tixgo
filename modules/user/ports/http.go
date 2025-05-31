package ports

import (
	"net/http"

	"tixgo/modules/user/app"
	"tixgo/modules/user/app/command"
	"tixgo/modules/user/app/query"
	"tixgo/shared/auth"
	"tixgo/shared/middleware"
	"tixgo/shared/response"
	"tixgo/shared/syserr"

	"github.com/gin-gonic/gin"
)

// HTTPHandler handles HTTP requests for user operations
type HTTPHandler struct {
	userService *app.UserService
	jwtService  *auth.JWTService
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(userService *app.UserService, jwtService *auth.JWTService) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// RegisterRoutes registers the user routes
func (h *HTTPHandler) RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", h.RegisterUser)
		userGroup.POST("/verify-otp", h.VerifyOTP)
		userGroup.POST("/login", h.LoginUser)

		userGroup.Use(middleware.RequireAuth(h.jwtService))
		userGroup.GET("/profile", h.GetUserProfile)
	}
}

// RegisterUser handles user registration
func (h *HTTPHandler) RegisterUser(c *gin.Context) {
	var req command.RegisterUserCommand
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	cmd := command.RegisterUserCommand{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserType:  req.UserType,
	}

	result, err := h.userService.RegisterUser(c.Request.Context(), cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewSimpleSuccessResponse(result).JSON(c, http.StatusCreated)
}

// VerifyOTP handles OTP verification
func (h *HTTPHandler) VerifyOTP(c *gin.Context) {
	var req command.VerifyOTPCommand
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	cmd := command.VerifyOTPCommand{
		Email: req.Email,
		OTP:   req.OTP,
	}

	result, err := h.userService.VerifyOTP(c.Request.Context(), cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewSimpleSuccessResponse(result).JSON(c, http.StatusOK)
}

// LoginUser handles user login
func (h *HTTPHandler) LoginUser(c *gin.Context) {
	var req command.LoginUserCommand
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	cmd := command.LoginUserCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.userService.LoginUser(c.Request.Context(), cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewSimpleSuccessResponse(result).JSON(c, http.StatusOK)
}

// GetUserProfile handles getting user profile
func (h *HTTPHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(syserr.New(syserr.UnauthorizedCode, "user not authenticated"))
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		c.Error(syserr.New(syserr.InternalCode, "invalid user ID"))
		return
	}

	query := query.GetUserProfileQuery{
		UserID: userIDInt64,
	}

	result, err := h.userService.GetUserProfile(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewSimpleSuccessResponse(result).JSON(c, http.StatusOK)
}
