package command

import (
	"context"

	"tixgo/modules/template/domain"

	"github.com/duongptryu/gox/syserr"
)

// UpdateTemplateCommand represents the command to update a template
type UpdateTemplateCommand struct {
	ID          int64    `json:"-"`
	Name        string   `json:"name"`
	Subject     string   `json:"subject"`
	Content     string   `json:"content"`
	Variables   []string `json:"variables"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
}

// UpdateTemplateResult represents the result of template update
type UpdateTemplateResult struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Subject     string                `json:"subject"`
	Type        domain.TemplateType   `json:"type"`
	Status      domain.TemplateStatus `json:"status"`
	Variables   []string              `json:"variables"`
	Description string                `json:"description"`
	UpdatedAt   string                `json:"updated_at"`
}

// UpdateTemplateHandler handles template updates
type UpdateTemplateHandler struct {
	templateRepo     domain.TemplateRepository
	templateRenderer domain.TemplateRenderer
}

// NewUpdateTemplateHandler creates a new update template handler
func NewUpdateTemplateHandler(templateRepo domain.TemplateRepository, templateRenderer domain.TemplateRenderer) *UpdateTemplateHandler {
	return &UpdateTemplateHandler{
		templateRepo:     templateRepo,
		templateRenderer: templateRenderer,
	}
}

// Handle executes the update template command
func (h *UpdateTemplateHandler) Handle(ctx context.Context, cmd UpdateTemplateCommand) error {
	// Get existing template
	template, err := h.templateRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		if err == domain.ErrTemplateNotFound {
			return domain.ErrTemplateNotFound
		}
		return syserr.Wrap(err, syserr.InternalCode, "failed to get template")
	}

	// Validate template content if provided
	if cmd.Content != "" {
		err = h.templateRenderer.ValidateTemplate(ctx, cmd.Content)
		if err != nil {
			return syserr.Wrap(err, syserr.InvalidArgumentCode, "template syntax validation failed")
		}
	}

	// Update template
	template.Update(cmd.Name, cmd.Subject, cmd.Content, cmd.Description, cmd.Variables)

	// Update status if provided
	if cmd.Status != "" {
		switch domain.TemplateStatus(cmd.Status) {
		case domain.TemplateStatusActive:
			template.Activate()
		case domain.TemplateStatusInactive:
			template.Deactivate()
		case domain.TemplateStatusDraft:
			template.Status = domain.TemplateStatusDraft
		default:
			return domain.ErrInvalidTemplateStatus
		}
	}

	// Save updated template
	err = h.templateRepo.Update(ctx, template)
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to update template")
	}

	return nil
}
