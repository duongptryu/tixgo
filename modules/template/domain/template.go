package domain

import (
	"time"

	"github.com/duongptryu/gox/syserr"
)

// TemplateType represents the type of template
type TemplateType string

const (
	TemplateTypeEmail TemplateType = "email"
	TemplateTypeSMS   TemplateType = "sms"
	TemplateTypePush  TemplateType = "push"
)

// TemplateStatus represents the status of template
type TemplateStatus string

const (
	TemplateStatusActive   TemplateStatus = "active"
	TemplateStatusInactive TemplateStatus = "inactive"
	TemplateStatusDraft    TemplateStatus = "draft"
)

// Template represents the template aggregate root
type Template struct {
	ID          int64
	Name        string
	Slug        string
	Subject     string
	Content     string
	Type        TemplateType
	Status      TemplateStatus
	Variables   []string
	Description string
	CreatedBy   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTemplate creates a new template
func NewTemplate(name, slug, subject, content string, templateType TemplateType, variables []string, description string, createdBy int64) (*Template, error) {
	if name == "" {
		return nil, syserr.New(syserr.InvalidArgumentCode, "template name is required")
	}
	if slug == "" {
		return nil, syserr.New(syserr.InvalidArgumentCode, "template slug is required")
	}
	if content == "" {
		return nil, syserr.New(syserr.InvalidArgumentCode, "template content is required")
	}
	if !IsValidTemplateType(string(templateType)) {
		return nil, syserr.New(syserr.InvalidArgumentCode, "invalid template type")
	}

	now := time.Now()
	return &Template{
		Name:        name,
		Slug:        slug,
		Subject:     subject,
		Content:     content,
		Type:        templateType,
		Status:      TemplateStatusDraft,
		Variables:   variables,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Activate sets the template status to active
func (t *Template) Activate() {
	t.Status = TemplateStatusActive
	t.UpdatedAt = time.Now()
}

// Deactivate sets the template status to inactive
func (t *Template) Deactivate() {
	t.Status = TemplateStatusInactive
	t.UpdatedAt = time.Now()
}

// Update updates the template content and metadata
func (t *Template) Update(name, subject, content, description string, variables []string) {
	if name != "" {
		t.Name = name
	}
	if subject != "" {
		t.Subject = subject
	}
	if content != "" {
		t.Content = content
	}
	if description != "" {
		t.Description = description
	}
	if variables != nil {
		t.Variables = variables
	}
	t.UpdatedAt = time.Now()
}

// IsActive checks if the template is active
func (t *Template) IsActive() bool {
	return t.Status == TemplateStatusActive
}

// IsValidTemplateType checks if the template type is valid
func IsValidTemplateType(templateType string) bool {
	switch TemplateType(templateType) {
	case TemplateTypeEmail, TemplateTypeSMS, TemplateTypePush:
		return true
	default:
		return false
	}
}
