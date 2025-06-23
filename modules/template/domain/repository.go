package domain

import (
	"context"

	"github.com/duongptryu/gox/pagination"
)

// TemplateRepository defines the interface for template persistence
type TemplateRepository interface {
	// Create creates a new template
	Create(ctx context.Context, template *Template) error

	// GetByID retrieves a template by ID
	GetByID(ctx context.Context, id int64) (*Template, error)

	// GetBySlug retrieves a template by slug
	GetBySlug(ctx context.Context, slug string) (*Template, error)

	// List retrieves templates with pagination and filters
	List(ctx context.Context, filters ListTemplateFilters, paging *pagination.Paging) ([]*Template, error)

	// Update updates an existing template
	Update(ctx context.Context, template *Template) error

	// Delete deletes a template by ID
	Delete(ctx context.Context, id int64) error
}

// TemplateRenderer defines the interface for template rendering
type TemplateRenderer interface {
	// Render renders a template with given variables
	Render(ctx context.Context, template *Template, variables map[string]interface{}) (*RenderedTemplate, error)

	// ValidateTemplate validates template syntax
	ValidateTemplate(ctx context.Context, content string) error
}

// ListTemplateFilters represents filters for listing templates
type ListTemplateFilters struct {
	Type      *TemplateType
	Status    *TemplateStatus
	CreatedBy *int64
	Search    string
}

// RenderedTemplate represents a rendered template result
type RenderedTemplate struct {
	Subject     string
	Content     string
	ContentType string
}
