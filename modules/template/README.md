# Template Module

The Template Module provides email/SMS/push notification template management and rendering capabilities for the TixGo platform. It follows Clean Architecture principles and supports dynamic template rendering with variables.

## Features

- **Template Management**: Create, update, delete, and list templates
- **Template Rendering**: Render templates with dynamic variables
- **Multiple Template Types**: Support for email, SMS, and push notification templates
- **Template Validation**: Validate template syntax before saving
- **Rich Template Functions**: Built-in helper functions for text manipulation
- **Status Management**: Draft, active, and inactive template states

## Architecture

The module follows Clean Architecture with these layers:

```
modules/template/
├── domain/          # Business logic and entities
├── app/            
│   ├── command/    # Write operations (create, update)
│   └── query/      # Read operations (get, list, render)
├── adapters/       # Infrastructure (database, template engine)
└── ports/          # HTTP handlers
```

## API Endpoints

### Public Endpoints
- `POST /api/templates/render` - Render a template with variables
- `GET /api/templates/by-slug/:slug` - Get template by slug

### Protected Endpoints (require authentication)
- `POST /api/templates` - Create a new template
- `GET /api/templates` - List templates with filters
- `GET /api/templates/:id` - Get template by ID
- `PUT /api/templates/:id` - Update template
- `DELETE /api/templates/:id` - Delete template

## Template Types

- **email**: HTML email templates with subject and content
- **sms**: Plain text SMS templates  
- **push**: Push notification templates

## Template Syntax

Templates use Go's `html/template` syntax with additional helper functions:

### Basic Variables
```html
<h1>Hello {{.Name}}!</h1>
<p>Your email is: {{.Email}}</p>
```

### Helper Functions
- `{{upper .Text}}` - Convert to uppercase
- `{{lower .Text}}` - Convert to lowercase  
- `{{title .Text}}` - Convert to title case
- `{{trim .Text}}` - Remove whitespace
- `{{default "fallback" .Value}}` - Use fallback if value is empty
- `{{contains .Text "substring"}}` - Check if text contains substring
- `{{replace .Text "old" "new"}}` - Replace text

### Conditional Logic
```html
{{if .ShowButton}}
<a href="{{.ButtonLink}}">{{.ButtonText}}</a>
{{end}}
```

### Loops
```html
{{range .Items}}
<li>{{.Name}}: {{.Price}}</li>
{{end}}
```

## Pagination

The module uses the [gox pagination library](https://github.com/duongptryu/gox/blob/main/pagination/pagination.go) for consistent pagination across the platform.

### Query Parameters for Listing:
- `page` - Page number (starts from 1)
- `limit` - Number of items per page (default: 20, max: 100)
- `type` - Filter by template type (email, sms, push)
- `status` - Filter by status (active, inactive, draft)
- `created_by` - Filter by creator user ID
- `search` - Search in name, description, or slug

### Pagination Response:
```json
{
  "page": 1,
  "limit": 20,
  "total": 85,
  "next_cursor": 21
}
```

The pagination object includes methods:
- `HasNext()` - Check if there's a next page
- `HasPrev()` - Check if there's a previous page
- `GetTotalPages()` - Calculate total pages
- `GetOffset()` - Calculate database offset

## Usage Examples

### 1. Create an Email Template

```bash
curl -X POST http://localhost:8080/api/templates \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Welcome Email",
    "slug": "welcome-email",
    "subject": "Welcome to {{.AppName}}!",
    "content": "<h1>Welcome {{.Name}}!</h1><p>Thanks for joining {{.AppName}}.</p>",
    "type": "email",
    "variables": ["AppName", "Name"],
    "description": "Welcome email for new users"
  }'
```

### 2. Render a Template

```bash
curl -X POST http://localhost:8080/api/templates/render \
  -H "Content-Type: application/json" \
  -d '{
    "template_slug": "welcome-email",
    "variables": {
      "AppName": "TixGo",
      "Name": "John Doe"
    }
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "subject": "Welcome to TixGo!",
    "content": "<h1>Welcome John Doe!</h1><p>Thanks for joining TixGo.</p>",
    "content_type": "text/html",
    "template_id": 1
  }
}
```

### 3. List Templates

```bash
curl "http://localhost:8080/api/templates?type=email&status=active&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:
```json
{
  "success": true,
  "data": {
    "templates": [...],
    "paging": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "next_cursor": 11
    }
  }
}
```

### 4. Get Template

```bash
curl http://localhost:8080/api/templates/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Email Template Examples

### OTP Verification Email
```html
<!DOCTYPE html>
<html>
<head>
    <title>Email Verification</title>
</head>
<body>
    <div style="max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif;">
        <h1>{{.AppName}} - Email Verification</h1>
        <p>Hello {{default "User" .Name}},</p>
        <p>Your verification code is: <strong>{{.OTP}}</strong></p>
        <p>This code will expire in {{default "10" .ExpiryMinutes}} minutes.</p>
        <p>Best regards,<br>The {{.AppName}} Team</p>
    </div>
</body>
</html>
```

### Password Reset Email
```html
<!DOCTYPE html>
<html>
<body>
    <div style="max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif;">
        <h1>Password Reset Request</h1>
        <p>Hello {{.Name}},</p>
        <p>You requested a password reset for your {{.AppName}} account.</p>
        <p><a href="{{.ResetLink}}" style="background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Reset Password</a></p>
        <p>This link will expire in {{.ExpiryHours}} hours.</p>
        <p>If you didn't request this, please ignore this email.</p>
    </div>
</body>
</html>
```

## Template Variables Best Practices

1. **Define Variables**: Always include a `variables` array when creating templates
2. **Use Defaults**: Use the `default` function for optional variables
3. **Validate Input**: Ensure all required variables are provided when rendering
4. **Escape HTML**: Use `safeHTML` function only for trusted content
5. **Test Templates**: Always test template rendering before activating

## Integration with Email Service

The template module can be integrated with email services:

```go
// Example: Send welcome email
templateRepo := adapters.NewTemplatePostgresRepository(db)
renderer := adapters.NewHTMLTemplateRenderer()
renderHandler := query.NewRenderTemplateHandler(templateRepo, renderer)

result, err := renderHandler.Handle(ctx, query.RenderTemplateQuery{
    TemplateSlug: &welcomeEmailSlug,
    Variables: map[string]interface{}{
        "Name": user.Name,
        "AppName": "TixGo",
    },
})

if err != nil {
    return err
}

// Send email using your email service
emailService.Send(user.Email, result.Subject, result.Content)
```

## Database Schema

```sql
CREATE TABLE templates (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    subject VARCHAR(500),
    content TEXT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('email', 'sms', 'push')),
    status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('active', 'inactive', 'draft')),
    variables TEXT[],
    description TEXT,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

## Testing

Run tests for the template module:

```bash
go test ./modules/template/...
```

The module includes comprehensive tests for:
- Template rendering with various scenarios
- Template validation
- Helper functions
- Error handling

## Error Handling

The module defines specific domain errors:
- `ErrTemplateNotFound` - Template doesn't exist
- `ErrTemplateAlreadyExists` - Template slug already in use
- `ErrInvalidTemplateType` - Invalid template type
- `ErrTemplateInactive` - Template is not active
- `ErrTemplateSyntaxError` - Template syntax is invalid

## Security Considerations

1. **Input Validation**: Always validate template content and variables
2. **XSS Prevention**: Use Go's html/template for automatic HTML escaping
3. **Access Control**: Only authenticated users can manage templates
4. **Template Isolation**: Each template is isolated during rendering 