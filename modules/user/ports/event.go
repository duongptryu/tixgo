package ports

import (
	"context"
	"fmt"
	"tixgo/components"
	"tixgo/modules/user/domain"

	userEvent "tixgo/modules/user/app/event"

	"github.com/duongptryu/gox/messaging"
)

const (
	EventUserRegistered = "user-registered"
)

func RegisterUserEventHandlers(eventBus messaging.Dispatcher, appCtx components.AppContext) {
	eventBus.RegisterEventHandler(EventUserRegistered, HandleEventUserRegistered(appCtx))
}

func HandleEventUserRegistered(appCtx components.AppContext) messaging.EventHandler {
	return func(ctx context.Context, event *any) error {
		userRegisteredEvent, ok := (*event).(*domain.EventUserRegistered)
		if !ok {
			return fmt.Errorf("invalid event type, expected *domain.EventUserRegistered")
		}

		biz := userEvent.NewSendMailOnUserRegistered(appCtx.GetCommandBus())

		err := biz.SendMailVerification(ctx, userRegisteredEvent)
		if err != nil {
			return err
		}

		return nil
	}
}
