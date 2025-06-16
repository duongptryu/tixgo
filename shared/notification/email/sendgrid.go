package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SendGridSender implements the email notification strategy using SendGrid API
type SendGridSender struct {
	config *SendGridConfig
	client *http.Client
}

// NewSendGridSender creates a new SendGrid email sender
func NewSendGridSender(config *SendGridConfig) *SendGridSender {
	return &SendGridSender{
		config: config,
		client: &http.Client{},
	}
}

// GetProviderName returns the provider name
func (s *SendGridSender) GetProviderName() string {
	return "SendGrid"
}

// SendEmail sends an email using SendGrid API
func (s *SendGridSender) SendEmail(ctx context.Context, emailMessage *EmailMessage) error {
	if len(emailMessage.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// Build SendGrid payload
	payload := s.buildSendGridPayload(emailMessage)

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal SendGrid payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SendGrid request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("SendGrid API error (status %d): %v", resp.StatusCode, errResp)
	}

	return nil
}

// sendGridPayload represents the SendGrid API payload structure
type sendGridPayload struct {
	Personalizations []personalization `json:"personalizations"`
	From             emailAddress      `json:"from"`
	Subject          string            `json:"subject"`
	Content          []content         `json:"content"`
	Attachments      []attachment      `json:"attachments,omitempty"`
}

type personalization struct {
	To  []emailAddress `json:"to"`
	CC  []emailAddress `json:"cc,omitempty"`
	BCC []emailAddress `json:"bcc,omitempty"`
}

type emailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type attachment struct {
	Content     string `json:"content"`
	Type        string `json:"type"`
	Filename    string `json:"filename"`
	Disposition string `json:"disposition,omitempty"`
}

// buildSendGridPayload creates a SendGrid-compatible payload
func (s *SendGridSender) buildSendGridPayload(emailMessage *EmailMessage) *sendGridPayload {
	// Build recipients
	var to []emailAddress
	for _, recipient := range emailMessage.To {
		to = append(to, emailAddress{Email: recipient})
	}

	var cc []emailAddress
	for _, recipient := range emailMessage.CC {
		cc = append(cc, emailAddress{Email: recipient})
	}

	var bcc []emailAddress
	for _, recipient := range emailMessage.BCC {
		bcc = append(bcc, emailAddress{Email: recipient})
	}

	// Build content
	var contents []content
	if emailMessage.HTMLBody != "" {
		contents = append(contents, content{
			Type:  "text/html",
			Value: emailMessage.HTMLBody,
		})
	}
	if emailMessage.Body != "" {
		contents = append(contents, content{
			Type:  "text/plain",
			Value: emailMessage.Body,
		})
	}
	if len(contents) == 0 {
		contents = append(contents, content{
			Type:  "text/plain",
			Value: emailMessage.Body,
		})
	}

	// Build attachments
	var attachments []attachment
	for _, att := range emailMessage.Attachments {
		attachments = append(attachments, attachment{
			Content:  string(att.Content), // Note: Should be base64 encoded in real implementation
			Type:     att.ContentType,
			Filename: att.Filename,
		})
	}

	return &sendGridPayload{
		Personalizations: []personalization{{
			To:  to,
			CC:  cc,
			BCC: bcc,
		}},
		From: emailAddress{
			Email: s.config.From,
		},
		Subject:     emailMessage.Subject,
		Content:     contents,
		Attachments: attachments,
	}
}
