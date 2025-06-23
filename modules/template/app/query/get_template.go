package query

import (
	"context"

	"tixgo/modules/template/domain"

	"github.com/duongptryu/gox/syserr"
)

// GetTemplateQuery represents the query to get a template
type GetTemplateQuery struct {
	ID   *int64  `json:"id"`
	Slug *string `json:"slug"`
}

// TemplateResult represents the template result
type TemplateResult struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Subject     string                `json:"subject"`
	Content     string                `json:"content"`
	Type        domain.TemplateType   `json:"type"`
	Status      domain.TemplateStatus `json:"status"`
	Variables   []string              `json:"variables"`
	Description string                `json:"description"`
	CreatedBy   int64                 `json:"created_by"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
}

// GetTemplateHandler handles getting template
type GetTemplateHandler struct {
	templateRepo domain.TemplateRepository
}

// NewGetTemplateHandler creates a new get template handler
func NewGetTemplateHandler(templateRepo domain.TemplateRepository) *GetTemplateHandler {
	return &GetTemplateHandler{
		templateRepo: templateRepo,
	}
}

// Handle executes the get template query
func (h *GetTemplateHandler) Handle(ctx context.Context, query GetTemplateQuery) (*TemplateResult, error) {
	var template *domain.Template
	var err error

	if query.ID != nil {
		template, err = h.templateRepo.GetByID(ctx, *query.ID)
	} else if query.Slug != nil {
		template, err = h.templateRepo.GetBySlug(ctx, *query.Slug)
	} else {
		return nil, syserr.New(syserr.InvalidArgumentCode, "either id or slug must be provided")
	}

	if err != nil {
		if err == domain.ErrTemplateNotFound {
			return nil, domain.ErrTemplateNotFound
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get template")
	}

	return &TemplateResult{
		ID:          template.ID,
		Name:        template.Name,
		Slug:        template.Slug,
		Subject:     template.Subject,
		Content:     template.Content,
		Type:        template.Type,
		Status:      template.Status,
		Variables:   template.Variables,
		Description: template.Description,
		CreatedBy:   template.CreatedBy,
		CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   template.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
