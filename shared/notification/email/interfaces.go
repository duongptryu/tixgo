package email

import "context"

// Strategy names for email providers
const (
	StrategyNameSMTP     = "smtp"
	StrategyNameSendGrid = "sendgrid"
)

// EmailSender defines the interface for email notification strategies
type EmailSender interface {
	SendEmail(ctx context.Context, emailMessage *EmailMessage) error
	GetProviderName() string
}

// EmailMessage represents an email message with email-specific fields
type EmailMessage struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	HTMLBody    string
	Attachments []Attachment
	Headers     map[string]string
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string
	ContentType string
	Content     []byte
}

// EmailResult represents the result of sending an email
type EmailResult struct {
	Success   bool
	MessageID string
	Error     string
	Metadata  map[string]string
	Provider  string
}

// SMTPConfig holds SMTP server configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

// SendGridConfig holds SendGrid API configuration
type SendGridConfig struct {
	APIKey string
	From   string
}
