package command

import (
	"context"

	"tixgo/modules/template/domain"

	"github.com/duongptryu/gox/syserr"
)

// CreateTemplateCommand represents the command to create a new template
type CreateTemplateCommand struct {
	Name        string   `json:"name" validate:"required"`
	Slug        string   `json:"slug" validate:"required"`
	Subject     string   `json:"subject"`
	Content     string   `json:"content" validate:"required"`
	Type        string   `json:"type" validate:"required"`
	Variables   []string `json:"variables"`
	Description string   `json:"description"`
	CreatedBy   int64    `json:"-"`
}

// CreateTemplateResult represents the result of template creation
type CreateTemplateResult struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Subject     string                `json:"subject"`
	Type        domain.TemplateType   `json:"type"`
	Status      domain.TemplateStatus `json:"status"`
	Variables   []string              `json:"variables"`
	Description string                `json:"description"`
	CreatedAt   string                `json:"created_at"`
}

// CreateTemplateHandler handles template creation
type CreateTemplateHandler struct {
	templateRepo     domain.TemplateRepository
	templateRenderer domain.TemplateRenderer
}

// NewCreateTemplateHandler creates a new create template handler
func NewCreateTemplateHandler(templateRepo domain.TemplateRepository, templateRenderer domain.TemplateRenderer) *CreateTemplateHandler {
	return &CreateTemplateHandler{
		templateRepo:     templateRepo,
		templateRenderer: templateRenderer,
	}
}

// Handle executes the create template command
func (h *CreateTemplateHandler) Handle(ctx context.Context, cmd CreateTemplateCommand) (*CreateTemplateResult, error) {
	// Validate template type
	if !domain.IsValidTemplateType(cmd.Type) {
		return nil, domain.ErrInvalidTemplateType
	}

	// Check if template with slug already exists
	existingTemplate, err := h.templateRepo.GetBySlug(ctx, cmd.Slug)
	if err != nil && err != domain.ErrTemplateNotFound {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to check existing template")
	}
	if existingTemplate != nil {
		return nil, domain.ErrTemplateAlreadyExists
	}

	// Validate template syntax
	err = h.templateRenderer.ValidateTemplate(ctx, cmd.Content)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InvalidArgumentCode, "template syntax validation failed")
	}

	// Create new template
	template, err := domain.NewTemplate(
		cmd.Name,
		cmd.Slug,
		cmd.Subject,
		cmd.Content,
		domain.TemplateType(cmd.Type),
		cmd.Variables,
		cmd.Description,
		cmd.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	// Save template
	err = h.templateRepo.Create(ctx, template)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to create template")
	}

	return &CreateTemplateResult{
		ID:          template.ID,
		Name:        template.Name,
		Slug:        template.Slug,
		Subject:     template.Subject,
		Type:        template.Type,
		Status:      template.Status,
		Variables:   template.Variables,
		Description: template.Description,
		CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
