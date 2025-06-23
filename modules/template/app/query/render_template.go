package query

import (
	"context"

	"tixgo/modules/template/domain"

	"github.com/duongptryu/gox/syserr"
)

// RenderTemplateQuery represents the query to render a template
type RenderTemplateQuery struct {
	TemplateID   *int64                 `json:"template_id"`
	TemplateSlug *string                `json:"template_slug"`
	Variables    map[string]interface{} `json:"variables"`
}

// RenderTemplateResult represents the result of template rendering
type RenderTemplateResult struct {
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	TemplateID  int64  `json:"template_id"`
}

// RenderTemplateHandler handles template rendering
type RenderTemplateHandler struct {
	templateRepo     domain.TemplateRepository
	templateRenderer domain.TemplateRenderer
}

// NewRenderTemplateHandler creates a new render template handler
func NewRenderTemplateHandler(templateRepo domain.TemplateRepository, templateRenderer domain.TemplateRenderer) *RenderTemplateHandler {
	return &RenderTemplateHandler{
		templateRepo:     templateRepo,
		templateRenderer: templateRenderer,
	}
}

// Handle executes the render template query
func (h *RenderTemplateHandler) Handle(ctx context.Context, query RenderTemplateQuery) (*RenderTemplateResult, error) {
	var template *domain.Template
	var err error

	// Get template by ID or slug
	if query.TemplateID != nil {
		template, err = h.templateRepo.GetByID(ctx, *query.TemplateID)
	} else if query.TemplateSlug != nil {
		template, err = h.templateRepo.GetBySlug(ctx, *query.TemplateSlug)
	} else {
		return nil, syserr.New(syserr.InvalidArgumentCode, "either template_id or template_slug must be provided")
	}

	if err != nil {
		if err == domain.ErrTemplateNotFound {
			return nil, domain.ErrTemplateNotFound
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get template")
	}

	// Check if template is active
	if !template.IsActive() {
		return nil, domain.ErrTemplateInactive
	}

	// Render template
	rendered, err := h.templateRenderer.Render(ctx, template, query.Variables)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to render template")
	}

	return &RenderTemplateResult{
		Subject:     rendered.Subject,
		Content:     rendered.Content,
		ContentType: rendered.ContentType,
		TemplateID:  template.ID,
	}, nil
}
