package config

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type AppConfig struct {
	App      App      `mapstructure:"app"`
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
	JWT      JWT      `mapstructure:"jwt"`
	Kafka    Kafka    `mapstructure:"kafka"`
}

type App struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment" validate:"required,oneof=dev stg prod"`
	DebugMode   bool   `mapstructure:"debug_mode" validate:"required"`
}

type Server struct {
	Host         string        `mapstructure:"host" validate:"required,hostname"`
	Port         int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" validate:"required,min=1s"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" validate:"required,min=1s"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" validate:"required,min=1s"`
}

type Database struct {
	Type          string        `mapstructure:"type" validate:"required,oneof=postgres mysql sqlite"`
	Host          string        `mapstructure:"host" validate:"required,hostname"`
	Port          int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	User          string        `mapstructure:"user" validate:"required,alphanum"`
	Password      string        `mapstructure:"password" validate:"required,alphanum"`
	Name          string        `mapstructure:"name" validate:"required,ascii"`
	SSLMode       string        `mapstructure:"ssl_mode" validate:"omitempty,oneof=disable prefer require verify-ca verify-full"`
	MaxOpenConns  int           `mapstructure:"max_open_conns" validate:"required,min=1"`
	MaxIdleConns  int           `mapstructure:"max_idle_conns" validate:"required,min=1"`
	MaxLifetime   time.Duration `mapstructure:"max_lifetime" validate:"required,min=1s"`
	MaxIdleTime   time.Duration `mapstructure:"max_idle_time" validate:"required,min=1s"`
	MigrationPath string        `mapstructure:"migration_path" validate:"required"`
}

type JWT struct {
	SecretKey          string        `mapstructure:"secret_key" validate:"required"`
	AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry" validate:"required,min=1s"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry" validate:"required,min=1s"`
}

// type Redis struct {
// 	Host     string `mapstructure:"host" validate:"required,hostname"`
// 	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
// 	Password string `mapstructure:"password"`
// 	DB       int    `mapstructure:"db"` // default 0
// }

type Kafka struct {
	Brokers []string `mapstructure:"brokers" validate:"required,min=1"`
}

func (c *AppConfig) Validate() error {
	return validator.New().Struct(c)
}
