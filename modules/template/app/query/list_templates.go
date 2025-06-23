package query

import (
	"context"

	"tixgo/modules/template/domain"

	"github.com/duongptryu/gox/pagination"
	"github.com/duongptryu/gox/syserr"
)

// FilterTemplatesQuery represents the filters for listing templates
type FilterTemplatesQuery struct {
	Type      *string `json:"type" form:"type"`
	Status    *string `json:"status" form:"status"`
	CreatedBy *int64  `json:"created_by" form:"created_by"`
	Search    string  `json:"search" form:"search"`
}

// ListTemplatesResult represents the result of template listing
type ListTemplatesResult struct {
	Templates []*TemplateListItem `json:"templates"`
	Paging    *pagination.Paging  `json:"paging"`
}

// TemplateListItem represents a template item in the list
type TemplateListItem struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Subject     string                `json:"subject"`
	Type        domain.TemplateType   `json:"type"`
	Status      domain.TemplateStatus `json:"status"`
	Description string                `json:"description"`
	CreatedBy   int64                 `json:"created_by"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
}

// ListTemplatesHandler handles listing templates
type ListTemplatesHandler struct {
	templateRepo domain.TemplateRepository
}

// NewListTemplatesHandler creates a new list templates handler
func NewListTemplatesHandler(templateRepo domain.TemplateRepository) *ListTemplatesHandler {
	return &ListTemplatesHandler{
		templateRepo: templateRepo,
	}
}

// Handle executes the list templates query
func (h *ListTemplatesHandler) Handle(ctx context.Context, filters FilterTemplatesQuery, paging *pagination.Paging) (*ListTemplatesResult, error) {
	// Ensure paging is not nil (should already be handled in HTTP layer)
	if paging == nil {
		paging = &pagination.Paging{}
		paging.Fulfill()
	}

	// Build domain filters from query filters
	domainFilters := domain.ListTemplateFilters{
		Search: filters.Search,
	}

	// Set type filter
	if filters.Type != nil && *filters.Type != "" {
		if !domain.IsValidTemplateType(*filters.Type) {
			return nil, domain.ErrInvalidTemplateType
		}
		templateType := domain.TemplateType(*filters.Type)
		domainFilters.Type = &templateType
	}

	// Set status filter
	if filters.Status != nil && *filters.Status != "" {
		templateStatus := domain.TemplateStatus(*filters.Status)
		switch templateStatus {
		case domain.TemplateStatusActive, domain.TemplateStatusInactive, domain.TemplateStatusDraft:
			domainFilters.Status = &templateStatus
		default:
			return nil, domain.ErrInvalidTemplateStatus
		}
	}

	// Set created by filter
	if filters.CreatedBy != nil {
		domainFilters.CreatedBy = filters.CreatedBy
	}

	// Get templates
	templates, err := h.templateRepo.List(ctx, domainFilters, paging)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to list templates")
	}

	// Convert to list items
	items := make([]*TemplateListItem, len(templates))
	for i, template := range templates {
		items[i] = &TemplateListItem{
			ID:          template.ID,
			Name:        template.Name,
			Slug:        template.Slug,
			Subject:     template.Subject,
			Type:        template.Type,
			Status:      template.Status,
			Description: template.Description,
			CreatedBy:   template.CreatedBy,
			CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   template.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &ListTemplatesResult{
		Templates: items,
		Paging:    paging,
	}, nil
}
