package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// TwilioSMSSender implements SMS notification strategy using Twilio API
type TwilioSMSSender struct {
	config *TwilioConfig
	client *http.Client
}

// NewTwilioSMSSender creates a new Twilio SMS sender
func NewTwilioSMSSender(config *TwilioConfig) *TwilioSMSSender {
	return &TwilioSMSSender{
		config: config,
		client: &http.Client{},
	}
}

// GetProviderName returns the provider name
func (t *TwilioSMSSender) GetProviderName() string {
	return "Twilio"
}

// SendSMS sends an SMS using Twilio API
func (t *TwilioSMSSender) SendSMS(ctx context.Context, smsMessage *SMSMessage) error {
	if len(smsMessage.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// Send SMS to each recipient
	for _, recipient := range smsMessage.To {
		err := t.sendSingleSMS(ctx, recipient, smsMessage.Message)
		if err != nil {
			return fmt.Errorf("failed to send SMS to %s: %w", recipient, err)
		}
	}

	return nil
}

// sendSingleSMS sends an SMS to a single recipient
func (t *TwilioSMSSender) sendSingleSMS(ctx context.Context, to, message string) error {
	// Twilio API endpoint
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.config.AccountSID)

	// Prepare form data
	data := url.Values{}
	data.Set("From", t.config.From)
	data.Set("To", to)
	data.Set("Body", message)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(t.config.AccountSID, t.config.AuthToken)

	// Send request
	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Twilio request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("Twilio API error (status %d): %v", resp.StatusCode, errResp)
	}

	return nil
}
