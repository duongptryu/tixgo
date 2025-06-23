package event

import (
	"github.com/duongptryu/gox/notification/mail"
)

type handleSendMailOtp struct {
	mailProvider mail.MailProvider
}

func NewHandleSendMailOtp(mailProvider mail.MailProvider) *handleSendMailOtp {
	return &handleSendMailOtp{
		mailProvider: mailProvider,
	}
}
