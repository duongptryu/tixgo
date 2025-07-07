package mail

import (
	"context"

	"github.com/duongptryu/gox/notification/mail"
)

type ConfigMail struct {
	OurMail string
	OurName string
}

type EventSendMailHandler struct {
	mailCfg      ConfigMail
	mailProvider mail.MailProvider
}

func NewEventSendMailHandler(mailProvider mail.MailProvider, cfgMail ConfigMail) *EventSendMailHandler {
	return &EventSendMailHandler{
		mailProvider: mailProvider,
		mailCfg:      cfgMail,
	}
}

func (h *EventSendMailHandler) Handle(ctx context.Context, event *EventSendMail) error {
	priority := mail.PriorityNormal
	if event.Priority != "" {
		priority = event.Priority
	}

	_, err := h.mailProvider.SendEmail(ctx, &mail.EmailMessage{
		From:     mail.EmailAddress{Email: h.mailCfg.OurMail, Name: h.mailCfg.OurName},
		To:       event.ToMail,
		CC:       event.CC,
		BCC:      event.BCC,
		Subject:  event.Subject,
		TextBody: event.TextBody,
		HTMLBody: event.HTMLBody,
		Priority: priority,
	})

	if err != nil {
		return err
	}

	return nil
}
