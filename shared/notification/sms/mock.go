package sms

import (
	"context"
	"fmt"
	"log"
)

// MockSMSSender implements a mock SMS notification strategy for testing
type MockSMSSender struct {
	config *MockSMSConfig
}

// NewMockSMSSender creates a new mock SMS sender
func NewMockSMSSender(config *MockSMSConfig) *MockSMSSender {
	return &MockSMSSender{
		config: config,
	}
}

// GetProviderName returns the provider name
func (m *MockSMSSender) GetProviderName() string {
	return "Mock"
}

// SendSMS sends a mock SMS message (logs to console)
func (m *MockSMSSender) SendSMS(ctx context.Context, smsMessage *SMSMessage) error {
	if len(smsMessage.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// Mock SMS sending - just log the message
	for _, recipient := range smsMessage.To {
		log.Printf("[MOCK SMS] From: %s, To: %s, Message: %s",
			m.config.From, recipient, smsMessage.Message)
	}

	return nil
}
