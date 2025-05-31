package components

import "github.com/jmoiron/sqlx"

type AppContext struct {
	db *sqlx.DB
}

func NewAppContext(db *sqlx.DB) *AppContext {
	return &AppContext{db: db}
}

func (c *AppContext) GetDB() *sqlx.DB {
	return c.db
}
