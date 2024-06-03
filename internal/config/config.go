package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPC       GRPC       `env-required:"true"`
	PostgreSQL PostgreSQL `env-required:"true"`
	Migrations Migrations `env-required:"true"`
}

type GRPC struct {
	Port    uint16        `env:"GRPC_PORT" env-default:"6969"`
	Timeout time.Duration `env:"GRPC_TIMEOUT" env-default:"5s"`
}

type PostgreSQL struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     uint16 `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_DBNAME" env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE" env-required:"true" env-default:"disable"`
}

type Migrations struct {
	Path string `env:"MIGRATIONS_PATH" env-required:"true"`
}

func MustLoadConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
