package adapters

import (
	"bytes"
	"context"
	"html/template"
	"strings"

	"tixgo/modules/template/domain"

	"github.com/duongptryu/gox/syserr"
)

// HTMLTemplateRenderer implements domain.TemplateRenderer using Go's html/template
type HTMLTemplateRenderer struct{}

// NewHTMLTemplateRenderer creates a new HTML template renderer
func NewHTMLTemplateRenderer() *HTMLTemplateRenderer {
	return &HTMLTemplateRenderer{}
}

// Render renders a template with given variables
func (r *HTMLTemplateRenderer) Render(ctx context.Context, tmpl *domain.Template, variables map[string]interface{}) (*domain.RenderedTemplate, error) {
	// Ensure variables is not nil
	if variables == nil {
		variables = make(map[string]interface{})
	}

	// Render subject
	renderedSubject, err := r.renderText(tmpl.Subject, variables)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to render subject")
	}

	// Render content
	renderedContent, err := r.renderHTML(tmpl.Content, variables)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to render content")
	}

	return &domain.RenderedTemplate{
		Subject:     renderedSubject,
		Content:     renderedContent,
		ContentType: "text/html",
	}, nil
}

// ValidateTemplate validates template syntax
func (r *HTMLTemplateRenderer) ValidateTemplate(ctx context.Context, content string) error {
	// Try to parse the template to check for syntax errors with helper functions
	tmpl := template.New("validation").Funcs(template.FuncMap{
		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"title":    strings.Title,
		"trim":     strings.TrimSpace,
		"contains": strings.Contains,
		"replace":  strings.ReplaceAll,
		"default": func(defaultValue interface{}, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
	})

	_, err := tmpl.Parse(content)
	if err != nil {
		return syserr.Wrap(err, syserr.InvalidArgumentCode, "template syntax error")
	}
	return nil
}

// renderText renders plain text template (for subjects)
func (r *HTMLTemplateRenderer) renderText(templateStr string, variables map[string]interface{}) (string, error) {
	if templateStr == "" {
		return "", nil
	}

	// Create template with helper functions (same as HTML template)
	tmpl := template.New("subject").Funcs(template.FuncMap{
		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"title":    strings.Title,
		"trim":     strings.TrimSpace,
		"contains": strings.Contains,
		"replace":  strings.ReplaceAll,
		"default": func(defaultValue interface{}, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
	})

	tmpl, err := tmpl.Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, variables)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

// renderHTML renders HTML template (for content)
func (r *HTMLTemplateRenderer) renderHTML(templateStr string, variables map[string]interface{}) (string, error) {
	if templateStr == "" {
		return "", nil
	}

	// Create template with helper functions
	tmpl := template.New("content").Funcs(template.FuncMap{
		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"title":    strings.Title,
		"trim":     strings.TrimSpace,
		"contains": strings.Contains,
		"replace":  strings.ReplaceAll,
		"default": func(defaultValue interface{}, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
	})

	tmpl, err := tmpl.Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, variables)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
