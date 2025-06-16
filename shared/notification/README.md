# Notification System

A clean and simple notification system using the Strategy pattern with direct provider configuration.

## Architecture

The notification system uses a direct provider approach:

- **EmailSender Interface**: Direct email notification interface
- **SMSSender Interface**: Direct SMS notification interface  
- **Provider Implementations**: SMTP, SendGrid for email; Twilio, Mock for SMS
- **Strategy Constants**: Predefined provider names to avoid string literals
- **Direct Configuration**: Configure providers once at application startup

## Key Benefits

- **Simple & Clean**: No complex managers - direct provider usage
- **Better Performance**: No map lookups or strategy resolution overhead
- **Clear Dependencies**: You know exactly which provider you're using
- **Easy Testing**: Mock specific provider interfaces directly
- **Dependency Injection Friendly**: Works perfectly with DI frameworks
- **Type Safe**: Provider constants prevent typos

## Running the Example

To see the notification system in action:

```bash
cd tixgo
go run examples/notification/main.go
```

This will demonstrate:
- Configuring providers at startup
- Sending email and SMS notifications directly
- Clean, simple API usage

For email testing with SMTP, you can set up MailHog:

```bash
# Install MailHog
go install github.com/mailhog/MailHog@latest

# Run MailHog (it will start on port 1025 for SMTP and 8025 for web UI)
MailHog

# Then run the example - emails will be captured by MailHog
go run examples/notification/main.go
```

## Usage

### Email Notifications

```go
import (
    "context"
    "tixgo/shared/notification/email"
)

// Configure provider once at startup
func setupEmailProvider() email.EmailSender {
    smtpConfig := &email.SMTPConfig{
        Host:     "smtp.gmail.com",
        Port:     587,
        Username: "your-email@gmail.com",
        Password: "your-app-password",
        From:     "your-email@gmail.com",
        UseTLS:   true,
    }
    return email.NewSMTPSender(smtpConfig)
}

// Use directly in your application
func SendWelcomeEmail(ctx context.Context, emailSender email.EmailSender, userEmail string) error {
    message := &email.EmailMessage{
        To:       []string{userEmail},
        Subject:  "Welcome to TixGo!",
        Body:     "Thank you for signing up!",
        HTMLBody: "<h1>Welcome!</h1><p>Thank you for signing up!</p>",
    }
    
    return emailSender.SendEmail(ctx, message)
}
```

### SMS Notifications

```go
import (
    "context"
    "tixgo/shared/notification/sms"
)

// Configure provider once at startup
func setupSMSProvider() sms.SMSSender {
    twilioConfig := &sms.TwilioConfig{
        AccountSID: "your-twilio-account-sid",
        AuthToken:  "your-twilio-auth-token",
        From:       "+1234567890",
    }
    return sms.NewTwilioSMSSender(twilioConfig)
}

// Use directly in your application
func SendVerificationCode(ctx context.Context, smsSender sms.SMSSender, phoneNumber, code string) error {
    message := &sms.SMSMessage{
        To:      []string{phoneNumber},
        Message: fmt.Sprintf("Your verification code is: %s", code),
        From:    "+1234567890",
    }
    
    return smsSender.SendSMS(ctx, message)
}
```

## Provider Configuration

### Application Structure

```go
// App dependencies configured at startup
type App struct {
    EmailSender email.EmailSender
    SMSSender   sms.SMSSender
}

func NewApp(config *Config) *App {
    return &App{
        EmailSender: createEmailSender(config.Email),
        SMSSender:   createSMSSender(config.SMS),
    }
}

func createEmailSender(config EmailConfig) email.EmailSender {
    switch config.Provider {
    case email.StrategyNameSMTP:
        return email.NewSMTPSender(&email.SMTPConfig{
            Host:     config.SMTP.Host,
            Port:     config.SMTP.Port,
            Username: config.SMTP.Username,
            Password: config.SMTP.Password,
            From:     config.SMTP.From,
            UseTLS:   config.SMTP.UseTLS,
        })
    case email.StrategyNameSendGrid:
        return email.NewSendGridSender(&email.SendGridConfig{
            APIKey: config.SendGrid.APIKey,
            From:   config.SendGrid.From,
        })
    default:
        panic("unknown email provider: " + config.Provider)
    }
}
```

### Configuration Example

```yaml
# config.yaml
email:
  provider: "smtp"  # Use email.StrategyNameSMTP
  smtp:
    host: "smtp.gmail.com"
    port: 587
    username: "your-email@gmail.com"
    password: "${SMTP_PASSWORD}"
    from: "noreply@tixgo.com"
    use_tls: true

sms:
  provider: "twilio"  # Use sms.StrategyNameTwilio
  twilio:
    account_sid: "${TWILIO_ACCOUNT_SID}"
    auth_token: "${TWILIO_AUTH_TOKEN}"
    from: "+1234567890"
```

## Provider Constants

### Email Provider Constants
```go
email.StrategyNameSMTP     // "smtp"
email.StrategyNameSendGrid // "sendgrid"
```

### SMS Provider Constants
```go
sms.StrategyNameMock   // "mock"  
sms.StrategyNameTwilio // "twilio"
```

## Available Providers

### Email Providers

#### SMTP Email
```go
smtpConfig := &email.SMTPConfig{
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your-email@gmail.com",
    Password: "your-app-password",
    From:     "your-email@gmail.com",
    UseTLS:   true,
}
emailSender := email.NewSMTPSender(smtpConfig)
```

#### SendGrid Email
```go
sendGridConfig := &email.SendGridConfig{
    APIKey: "your-sendgrid-api-key",
    From:   "noreply@yourdomain.com",
}
emailSender := email.NewSendGridSender(sendGridConfig)
```

### SMS Providers

#### Mock SMS (for testing)
```go
mockConfig := &sms.MockSMSConfig{
    From: "+1234567890",
}
smsSender := sms.NewMockSMSSender(mockConfig)
```

#### Twilio SMS
```go
twilioConfig := &sms.TwilioConfig{
    AccountSID: "your-twilio-account-sid",
    AuthToken:  "your-twilio-auth-token",
    From:       "+1234567890",
}
smsSender := sms.NewTwilioSMSSender(twilioConfig)
```

## Dependency Injection

Works perfectly with dependency injection frameworks:

```go
// Wire/Dig example
func ProvideEmailSender(config *Config) email.EmailSender {
    switch config.Email.Provider {
    case email.StrategyNameSMTP:
        return email.NewSMTPSender(config.Email.SMTP)
    case email.StrategyNameSendGrid:
        return email.NewSendGridSender(config.Email.SendGrid)
    default:
        panic("unknown email provider")
    }
}

// Usage in handlers
type UserHandler struct {
    emailSender email.EmailSender
    smsSender   sms.SMSSender
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    // ... registration logic ...
    
    // Send welcome email
    err := h.emailSender.SendEmail(r.Context(), welcomeMessage)
    if err != nil {
        log.Printf("Failed to send welcome email: %v", err)
    }
}
```

## Fallback Implementation

If you need fallback logic, implement it at the service level:

```go
type NotificationService struct {
    primaryEmailSender   email.EmailSender
    fallbackEmailSender email.EmailSender
}

func (s *NotificationService) SendEmail(ctx context.Context, message *email.EmailMessage) error {
    // Try primary provider
    if err := s.primaryEmailSender.SendEmail(ctx, message); err != nil {
        log.Printf("Primary email failed: %v", err)
        
        // Fallback to secondary provider
        if s.fallbackEmailSender != nil {
            return s.fallbackEmailSender.SendEmail(ctx, message)
        }
        return err
    }
    return nil
}
```

## Email Features

- **HTML and Text**: Support for both HTML and plain text content
- **Attachments**: File attachment support  
- **CC/BCC**: Carbon copy and blind carbon copy
- **Custom Headers**: Add custom email headers
- **TLS Support**: Secure email transmission

## Adding New Providers

### Adding a New Email Provider

1. Add provider constant to `interfaces.go`:

```go
const (
    StrategyNameSMTP       = "smtp"
    StrategyNameSendGrid   = "sendgrid"
    StrategyNameMyProvider = "myprovider" // Add your constant
)
```

2. Implement the `EmailSender` interface:

```go
type MyEmailProvider struct {
    config *MyConfig
}

func (m *MyEmailProvider) SendEmail(ctx context.Context, message *email.EmailMessage) error {
    // Your implementation
    return nil
}

func (m *MyEmailProvider) GetProviderName() string {
    return "MyProvider"
}
```

3. Create constructor function:

```go
func NewMyEmailProvider(config *MyConfig) *MyEmailProvider {
    return &MyEmailProvider{config: config}
}
```

### Adding a New SMS Provider

Follow the same pattern for SMS providers in the `sms` package.

## Testing

### Unit Testing

```go
func TestUserHandler_SendWelcomeEmail(t *testing.T) {
    // Create mock email sender
    mockEmailSender := &MockEmailSender{}
    
    handler := &UserHandler{
        emailSender: mockEmailSender,
    }
    
    // Test your handler
    err := handler.SendWelcomeEmail(context.Background(), "test@example.com")
    assert.NoError(t, err)
    assert.True(t, mockEmailSender.SendEmailCalled)
}

type MockEmailSender struct {
    SendEmailCalled bool
}

func (m *MockEmailSender) SendEmail(ctx context.Context, message *email.EmailMessage) error {
    m.SendEmailCalled = true
    return nil
}

func (m *MockEmailSender) GetProviderName() string {
    return "Mock"
}
```

### Integration Testing

```go
func TestEmailIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Use real provider with test configuration
    emailSender := email.NewSMTPSender(&email.SMTPConfig{
        Host: "localhost",
        Port: 1025, // MailHog
        From: "test@example.com",
    })
    
    message := &email.EmailMessage{
        To:      []string{"integration-test@example.com"},
        Subject: "Integration Test",
        Body:    "This is an integration test",
    }
    
    err := emailSender.SendEmail(context.Background(), message)
    assert.NoError(t, err)
}
```

## Error Handling

```go
func SendNotification(ctx context.Context, emailSender email.EmailSender, userEmail string) error {
    message := &email.EmailMessage{
        To:      []string{userEmail},
        Subject: "Notification",
        Body:    "You have a new notification",
    }
    
    if err := emailSender.SendEmail(ctx, message); err != nil {
        // Log error with context
        log.Printf("Failed to send notification to %s: %v", userEmail, err)
        
        // Handle specific error types if needed
        if isRateLimitError(err) {
            return fmt.Errorf("rate limited: %w", err)
        }
        
        return fmt.Errorf("notification failed: %w", err)
    }
    
    return nil
}
```

## Project Structure

```
shared/notification/
├── email/
│   ├── interfaces.go     # EmailSender interface, types, and constants
│   ├── smtp.go          # SMTP email implementation
│   └── sendgrid.go      # SendGrid email implementation
└── sms/
    ├── interfaces.go     # SMSSender interface, types, and constants
    ├── mock.go          # Mock SMS implementation
    └── twilio.go        # Twilio SMS implementation

examples/notification/
└── main.go              # Complete usage example
```

## Architecture Benefits

- **Simplicity**: Direct provider usage without management overhead
- **Performance**: No strategy resolution or map lookups
- **Clarity**: Clear provider dependencies in your application
- **Flexibility**: Easy to switch providers or implement fallbacks
- **Testing**: Simple mocking and testing strategies
- **Maintainability**: Less code to maintain and debug

This approach provides a clean, efficient, and maintainable notification system that's perfect for most real-world applications! 