package event

import (
	"context"
	"tixgo/modules/user/app/command"
	"tixgo/modules/user/domain"

	"github.com/duongptryu/gox/messaging"
)

type sendMailOnUserRegistered struct {
	commandBus messaging.CommandBus
}

func NewSendMailOnUserRegistered(commandBus messaging.CommandBus) *sendMailOnUserRegistered {
	return &sendMailOnUserRegistered{
		commandBus: commandBus,
	}
}

func (h *sendMailOnUserRegistered) SendMailVerification(ctx context.Context, event *domain.EventUserRegistered) error {
	sendMailVerificationCmd := &command.SendOTPVerifyMailCommand{
		Mail: event.Email,
	}

	return h.commandBus.PublishCommand(ctx, sendMailVerificationCmd)
}
