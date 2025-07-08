package ports

import (
	"context"
	"tixgo/components"

	templateAdapters "tixgo/modules/template/adapters"
	"tixgo/modules/user/adapters"
	"tixgo/modules/user/app/command"
	userEvent "tixgo/modules/user/app/event"
	"tixgo/modules/user/domain"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/duongptryu/gox/messaging"
)

const (
	EventUserRegistered      = "events.EventUserRegistered"
	CommandSendOTPVerifyMail = "commands.SendOTPVerifyMail"
)

type UserMessagingHandlers struct {
	dispatcher messaging.Dispatcher
	appCtx     components.AppContext
}

func NewUserMessagingHandlers(dispatcher messaging.Dispatcher, appCtx components.AppContext) *UserMessagingHandlers {
	return &UserMessagingHandlers{
		dispatcher: dispatcher,
		appCtx:     appCtx,
	}
}

func (h *UserMessagingHandlers) RegisterUserMessagingHandlers() {
	eventProcessor := h.dispatcher.GetEventProcessor()
	eventProcessor.AddHandler(cqrs.NewEventHandler(EventUserRegistered, h.HandleEventUserRegistered))

	commandProcessor := h.dispatcher.GetCommandProcessor()
	commandProcessor.AddHandler(cqrs.NewCommandHandler(CommandSendOTPVerifyMail, h.HandleCommandSendOTPVerifyMail))
}

func (h *UserMessagingHandlers) HandleEventUserRegistered(ctx context.Context, event *domain.EventUserRegistered) error {
	biz := userEvent.NewSendMailOnUserRegistered(h.appCtx.GetCommandBus())

	err := biz.SendMailVerification(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (h *UserMessagingHandlers) HandleCommandSendOTPVerifyMail(ctx context.Context, cmd *command.SendOTPVerifyMailCommand) error {
	otpStore := adapters.NewInMemoryOTPStore()
	templateRepo := templateAdapters.NewTemplatePostgresRepository(h.appCtx.GetDB())
	templateRenderer := templateAdapters.NewHTMLTemplateRenderer()
	biz := command.NewSendOTPVerifyMailHandler(otpStore, templateRepo, templateRenderer, h.appCtx.GetEventBus())

	err := biz.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
