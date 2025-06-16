package sms

import "context"

// Strategy names for SMS providers
const (
	StrategyNameMock   = "mock"
	StrategyNameTwilio = "twilio"
)

// SMSSender defines the interface for SMS notification strategies
type SMSSender interface {
	SendSMS(ctx context.Context, smsMessage *SMSMessage) error
	GetProviderName() string
}

// SMSMessage represents an SMS message
type SMSMessage struct {
	To      []string
	Message string
	From    string
	Data    map[string]string
}

// SMSResult represents the result of sending an SMS
type SMSResult struct {
	Success   bool
	MessageID string
	Error     string
	Metadata  map[string]string
	Provider  string
}

// TwilioConfig holds Twilio SMS configuration
type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	From       string
}

// MockSMSConfig holds configuration for mock SMS implementation
type MockSMSConfig struct {
	From string
}
