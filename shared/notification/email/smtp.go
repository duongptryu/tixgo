package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPSender implements the email notification strategy using SMTP
type SMTPSender struct {
	config *SMTPConfig
}

// NewSMTPSender creates a new SMTP email sender
func NewSMTPSender(config *SMTPConfig) *SMTPSender {
	return &SMTPSender{
		config: config,
	}
}

// GetProviderName returns the provider name
func (s *SMTPSender) GetProviderName() string {
	return "SMTP"
}

// SendEmail sends an email using SMTP
func (s *SMTPSender) SendEmail(ctx context.Context, emailMessage *EmailMessage) error {
	if len(emailMessage.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// Create the email message
	msg := s.buildMessage(emailMessage)

	// Setup SMTP authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// Server address
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Recipients list
	recipients := append(emailMessage.To, emailMessage.CC...)
	recipients = append(recipients, emailMessage.BCC...)

	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, s.config.From, recipients, msg)
	}

	return smtp.SendMail(addr, auth, s.config.From, recipients, []byte(msg))
}

// sendWithTLS sends email with TLS encryption
func (s *SMTPSender) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg string) error {
	// Create TLS connection
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return fmt.Errorf("failed to establish TLS connection: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Set sender
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// buildMessage constructs the email message
func (s *SMTPSender) buildMessage(emailMessage *EmailMessage) string {
	var msg strings.Builder

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", s.config.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(emailMessage.To, ",")))

	if len(emailMessage.CC) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(emailMessage.CC, ",")))
	}

	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", emailMessage.Subject))

	// Custom headers
	for key, value := range emailMessage.Headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Content type
	if emailMessage.HTMLBody != "" {
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	msg.WriteString("\r\n")

	// Body
	if emailMessage.HTMLBody != "" {
		msg.WriteString(emailMessage.HTMLBody)
	} else {
		msg.WriteString(emailMessage.Body)
	}

	return msg.String()
}
