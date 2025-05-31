package user

import (
	"tixgo/config"
	"tixgo/modules/user/adapters"
	"tixgo/modules/user/app"
	"tixgo/modules/user/ports"
	"tixgo/shared/auth"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Module represents the user module
type Module struct {
	HTTPHandler *ports.HTTPHandler
	UserService *app.UserService
}

// NewModule creates a new user module with all dependencies wired
func NewModule(db *sqlx.DB, jwtConfig config.JWT) (*Module, error) {
	// Create infrastructure adapters
	userRepo := adapters.NewUserPostgresRepository(db)
	otpStore := adapters.NewInMemoryOTPStore()

	// Create JWT service
	jwtService := auth.NewJWTService(
		jwtConfig.SecretKey,
		jwtConfig.AccessTokenExpiry,
		jwtConfig.RefreshTokenExpiry,
	)

	// Create application service
	userService := app.NewUserService(userRepo, otpStore, jwtService)

	// Create HTTP handler
	httpHandler := ports.NewHTTPHandler(userService, jwtService)

	return &Module{
		HTTPHandler: httpHandler,
		UserService: userService,
	}, nil
}

// RegisterRoutes registers the module's HTTP routes
func (m *Module) RegisterRoutes(router *gin.Engine) {
	m.HTTPHandler.RegisterRoutes(router)
}
