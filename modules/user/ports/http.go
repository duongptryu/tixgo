package ports

import (
	"net/http"

	"tixgo/components"
	"tixgo/modules/user/adapters"
	"tixgo/modules/user/app/command"
	"tixgo/modules/user/app/query"

	"github.com/duongptryu/gox/context"
	"github.com/duongptryu/gox/response"
	"github.com/duongptryu/gox/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, appCtx components.AppContext) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", RegisterUser(appCtx))
		userGroup.POST("/verify-otp", VerifyOTP(appCtx))
		userGroup.POST("/login", LoginUser(appCtx))

		userGroup.Use(middleware.RequireAuth(appCtx.GetJWTService()))
		userGroup.GET("/profile", GetUserProfile(appCtx))
	}
}

func RegisterUser(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req command.RegisterUserCommand
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			return
		}

		userRepo := adapters.NewUserPostgresRepository(appCtx.GetDB())
		tempUserStore := adapters.NewInMemoryTempUserStore()
		otpStore := adapters.NewInMemoryOTPStore()

		biz := command.NewRegisterUserHandler(userRepo, tempUserStore, otpStore, appCtx.GetEventBus())

		result, err := biz.Handle(c.Request.Context(), &req)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, response.NewSimpleSuccessResponse(result))
	}
}

func VerifyOTP(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req command.VerifyOTPCommand
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			return
		}

		userRepo := adapters.NewUserPostgresRepository(appCtx.GetDB())
		tempUserStore := adapters.NewInMemoryTempUserStore()
		otpStore := adapters.NewInMemoryOTPStore()

		biz := command.NewVerifyOTPHandler(userRepo, tempUserStore, otpStore)

		result, err := biz.Handle(c.Request.Context(), &req)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.NewSimpleSuccessResponse(result))
	}
}

func LoginUser(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req command.LoginUserCommand
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			return
		}

		userRepo := adapters.NewUserPostgresRepository(appCtx.GetDB())

		biz := command.NewLoginUserHandler(userRepo, appCtx.GetJWTService())

		result, err := biz.Handle(c.Request.Context(), &req)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.NewSimpleSuccessResponse(result))
	}
}

func GetUserProfile(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInt64, err := context.GetUserIDFromContextAsInt64(c.Request.Context())
		if err != nil {
			c.Error(err)
			return
		}

		userRepo := adapters.NewUserPostgresRepository(appCtx.GetDB())
		biz := query.NewGetUserProfileHandler(userRepo)

		result, err := biz.Handle(c.Request.Context(), &query.GetUserProfileQuery{
			UserID: userIDInt64,
		})
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.NewSimpleSuccessResponse(result))
	}
}
