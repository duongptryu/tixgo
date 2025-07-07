package components

import (
	"github.com/duongptryu/gox/auth"
	"github.com/duongptryu/gox/messaging"

	"github.com/jmoiron/sqlx"
)

type AppContext interface {
	GetDB() *sqlx.DB
	GetJWTService() *auth.JWTService
	GetCommandBus() messaging.CommandBus
	GetEventBus() messaging.EventBus
	GetDispatcher() messaging.Dispatcher
}

type appCtx struct {
	db         *sqlx.DB
	jwtService *auth.JWTService
	commandBus messaging.CommandBus
	eventBus   messaging.EventBus
	dispatcher messaging.Dispatcher
}

func NewAppContext(db *sqlx.DB, jwtService *auth.JWTService, commandBus messaging.CommandBus, eventBus messaging.EventBus, dispatcher messaging.Dispatcher) AppContext {
	return &appCtx{db: db, jwtService: jwtService, commandBus: commandBus, eventBus: eventBus, dispatcher: dispatcher}
}

func (c *appCtx) GetDB() *sqlx.DB {
	return c.db
}

func (c *appCtx) GetJWTService() *auth.JWTService {
	return c.jwtService
}

func (c *appCtx) GetCommandBus() messaging.CommandBus {
	return c.commandBus
}

func (c *appCtx) GetEventBus() messaging.EventBus {
	return c.eventBus
}

func (c *appCtx) GetDispatcher() messaging.Dispatcher {
	return c.dispatcher
}
