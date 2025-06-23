package adapters

import (
	"context"
	"testing"

	"tixgo/modules/template/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTMLTemplateRenderer_Render(t *testing.T) {
	renderer := NewHTMLTemplateRenderer()
	ctx := context.Background()

	tests := []struct {
		name      string
		template  *domain.Template
		variables map[string]interface{}
		expected  *domain.RenderedTemplate
		wantErr   bool
	}{
		{
			name: "simple email template",
			template: &domain.Template{
				Subject: "Welcome {{.Name}}!",
				Content: "<h1>Hello {{.Name}}</h1><p>Welcome to our platform!</p>",
			},
			variables: map[string]interface{}{
				"Name": "John Doe",
			},
			expected: &domain.RenderedTemplate{
				Subject:     "Welcome John Doe!",
				Content:     "<h1>Hello John Doe</h1><p>Welcome to our platform!</p>",
				ContentType: "text/html",
			},
			wantErr: false,
		},
		{
			name: "template with helper functions",
			template: &domain.Template{
				Subject: "Welcome {{upper .Name}}!",
				Content: `<h1>Hello {{title .Name}}</h1><p>Your email: {{lower .Email}}</p>`,
			},
			variables: map[string]interface{}{
				"Name":  "john doe",
				"Email": "JOHN@EXAMPLE.COM",
			},
			expected: &domain.RenderedTemplate{
				Subject:     "Welcome JOHN DOE!",
				Content:     "<h1>Hello John Doe</h1><p>Your email: john@example.com</p>",
				ContentType: "text/html",
			},
			wantErr: false,
		},
		{
			name: "template with default function",
			template: &domain.Template{
				Subject: "Hello {{default \"User\" .Name}}",
				Content: `<p>Phone: {{default "Not provided" .Phone}}</p>`,
			},
			variables: map[string]interface{}{
				"Name": "John",
			},
			expected: &domain.RenderedTemplate{
				Subject:     "Hello John",
				Content:     "<p>Phone: Not provided</p>",
				ContentType: "text/html",
			},
			wantErr: false,
		},
		{
			name: "empty variables map",
			template: &domain.Template{
				Subject: "Hello World",
				Content: "<p>Static content</p>",
			},
			variables: nil,
			expected: &domain.RenderedTemplate{
				Subject:     "Hello World",
				Content:     "<p>Static content</p>",
				ContentType: "text/html",
			},
			wantErr: false,
		},
		{
			name: "template with missing variable",
			template: &domain.Template{
				Subject: "Hello {{.Name}}",
				Content: "<p>Hello {{.Name}}</p>",
			},
			variables: map[string]interface{}{},
			expected: &domain.RenderedTemplate{
				Subject:     "Hello",
				Content:     "<p>Hello </p>",
				ContentType: "text/html",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderer.Render(ctx, tt.template, tt.variables)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.Subject, result.Subject)
			assert.Equal(t, tt.expected.Content, result.Content)
			assert.Equal(t, tt.expected.ContentType, result.ContentType)
		})
	}
}

func TestHTMLTemplateRenderer_ValidateTemplate(t *testing.T) {
	renderer := NewHTMLTemplateRenderer()
	ctx := context.Background()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid template",
			content: "<h1>Hello {{.Name}}</h1>",
			wantErr: false,
		},
		{
			name:    "valid template with functions",
			content: "<h1>Hello {{upper .Name}}</h1>",
			wantErr: false,
		},
		{
			name:    "invalid template syntax",
			content: "<h1>Hello {{.Name</h1>", // missing closing }}
			wantErr: true,
		},
		{
			name:    "empty template",
			content: "",
			wantErr: false,
		},
		{
			name:    "template with range",
			content: "{{range .Items}}<p>{{.}}</p>{{end}}",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := renderer.ValidateTemplate(ctx, tt.content)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHTMLTemplateRenderer_RenderComplexTemplate(t *testing.T) {
	renderer := NewHTMLTemplateRenderer()
	ctx := context.Background()

	template := &domain.Template{
		Subject: "OTP Verification - {{.AppName}}",
		Content: `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Subject}}</title>
</head>
<body>
    <div style="max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif;">
        <h1>{{title .AppName}} - Email Verification</h1>
        <p>Hello {{default "User" .Name}},</p>
        <p>Your OTP code is: <strong>{{.OTP}}</strong></p>
        <p>This code will expire in {{default "10" .ExpiryMinutes}} minutes.</p>
        {{if .LoginLink}}
        <p><a href="{{.LoginLink}}">Click here to login</a></p>
        {{end}}
        <p>Best regards,<br>The {{.AppName}} Team</p>
    </div>
</body>
</html>`,
	}

	variables := map[string]interface{}{
		"AppName":       "tixgo",
		"Name":          "John Doe",
		"OTP":           "123456",
		"ExpiryMinutes": "15",
		"LoginLink":     "https://app.tixgo.com/login",
	}

	result, err := renderer.Render(ctx, template, variables)

	require.NoError(t, err)
	assert.Equal(t, "OTP Verification - tixgo", result.Subject)
	assert.Contains(t, result.Content, "Tixgo - Email Verification")
	assert.Contains(t, result.Content, "Hello John Doe")
	assert.Contains(t, result.Content, "Your OTP code is: <strong>123456</strong>")
	assert.Contains(t, result.Content, "expire in 15 minutes")
	assert.Contains(t, result.Content, `<a href="https://app.tixgo.com/login">Click here to login</a>`)
	assert.Equal(t, "text/html", result.ContentType)
}
