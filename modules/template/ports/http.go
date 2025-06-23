package ports

import (
	"net/http"
	"strconv"

	"tixgo/components"
	"tixgo/modules/template/adapters"
	"tixgo/modules/template/app/command"
	"tixgo/modules/template/app/query"

	"github.com/duongptryu/gox/context"
	"github.com/duongptryu/gox/pagination"
	"github.com/duongptryu/gox/response"
	"github.com/duongptryu/gox/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTemplateRoutes(router *gin.RouterGroup, appCtx components.AppContext) {
	templateGroup := router.Group("/templates")
	{
		// Public endpoints for rendering templates
		templateGroup.POST("/render", RenderTemplate(appCtx))
		templateGroup.GET("/by-slug/:slug", GetTemplateBySlug(appCtx))

		// Protected endpoints requiring authentication
		templateGroup.Use(middleware.RequireAuth(appCtx.GetJWTService()))
		templateGroup.POST("", CreateTemplate(appCtx))
		templateGroup.GET("", ListTemplates(appCtx))
		templateGroup.GET("/:id", GetTemplate(appCtx))
		templateGroup.PUT("/:id", UpdateTemplate(appCtx))
		templateGroup.DELETE("/:id", DeleteTemplate(appCtx))
	}
}

func CreateTemplate(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req command.CreateTemplateCommand
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			return
		}

		// Get user ID from context
		userID, err := context.GetUserIDFromContextAsInt64(c.Request.Context())
		if err != nil {
			c.Error(err)
			return
		}
		req.CreatedBy = userID

		templateRepo := adapters.NewTemplatePostgresRepository(appCtx.GetDB())
		templateRenderer := adapters.NewHTMLTemplateRenderer()

		handler := command.NewCreateTemplateHandler(templateRepo, templateRenderer)

		result, err := handler.Handle(c.Request.Context(), req)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, response.NewSimpleSuccessResponse(result))
	}
}

func UpdateTemplate(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req command.UpdateTemplateCommand
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			return
		}

		// Get template ID from URL parameter
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.Error(err)
			return
		}
		req.ID = id

		templateRepo := adapters.NewTemplatePostgresRepository(appCtx.GetDB())
		templateRenderer := adapters.NewHTMLTemplateRenderer()

		handler := command.NewUpdateTemplateHandler(templateRepo, templateRenderer)

		result, err := handler.Handle(c.Request.Context(), req)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.NewSimpleSuccessResponse(result))
	}
}

func GetTemplate(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get template ID from URL parameter
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.Error(err)
			return
		}

		templateRepo := adapters.NewTemplatePostgresRepository(appCtx.GetDB())
		handler := query.NewGetTemplateHandler(templateRepo)

		result, err := handler.Handle(c.Request.Context(), query.GetTemplateQuery{
			ID: &id,
		})
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.NewSimpleSuccessResponse(result))
	}
}

func GetTemplateBySlug(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		templateRepo := adapters.NewTemplatePostgresRepository(appCtx.GetDB())
		handler := query.NewGetTemplateHandler(templateRepo)

		result, err := handler.Handle(c.Request.Context(), query.GetTemplateQuery{
			Slug: &slug,
		})
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.NewSimpleSuccessResponse(result))
	}
}

func ListTemplates(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind filters separately (ShouldBind is more forgiving for optional parameters)
		var filters query.FilterTemplatesQuery
		if err := c.ShouldBind(&filters); err != nil {
			c.Error(err)
			return
		}

		// Bind paging separately
		var paging pagination.Paging
		if err := c.ShouldBind(&paging); err != nil {
			c.Error(err)
			return
		}

		// Apply pagination defaults in HTTP layer
		paging.Fulfill()

		templateRepo := adapters.NewTemplatePostgresRepository(appCtx.GetDB())
		handler := query.NewListTemplatesHandler(templateRepo)

		result, err := handler.Handle(c.Request.Context(), filters, &paging)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.NewSimpleSuccessResponse(result))
	}
}

func RenderTemplate(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req query.RenderTemplateQuery
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			return
		}

		templateRepo := adapters.NewTemplatePostgresRepository(appCtx.GetDB())
		templateRenderer := adapters.NewHTMLTemplateRenderer()

		handler := query.NewRenderTemplateHandler(templateRepo, templateRenderer)

		result, err := handler.Handle(c.Request.Context(), req)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.NewSimpleSuccessResponse(result))
	}
}

func DeleteTemplate(appCtx components.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get template ID from URL parameter
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.Error(err)
			return
		}

		templateRepo := adapters.NewTemplatePostgresRepository(appCtx.GetDB())

		err = templateRepo.Delete(c.Request.Context(), id)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.NewSimpleSuccessResponse(map[string]string{
			"message": "Template deleted successfully",
		}))
	}
}
