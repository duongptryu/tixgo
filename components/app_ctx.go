package components

import (
	"github.com/duongptryu/gox/auth"

	"github.com/jmoiron/sqlx"
)

type AppContext interface {
	GetDB() *sqlx.DB
	GetJWTService() *auth.JWTService
}

type appCtx struct {
	db         *sqlx.DB
	jwtService *auth.JWTService
}

func NewAppContext(db *sqlx.DB, jwtService *auth.JWTService) AppContext {
	return &appCtx{db: db, jwtService: jwtService}
}

func (c *appCtx) GetDB() *sqlx.DB {
	return c.db
}

func (c *appCtx) GetJWTService() *auth.JWTService {
	return c.jwtService
}
