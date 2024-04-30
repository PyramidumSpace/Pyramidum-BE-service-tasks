package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	GRPC       GRPC       `yaml:"grpc" env-required:"true"`
	PostgreSQL PostgreSQL `yaml:"postgresql" env-required:"true"`
	Migrations Migrations `yaml:"migrations" env-required:"true"`
}

type GRPC struct {
	Port    uint16        `yaml:"port" env-default:"6969"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type PostgreSQL struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     uint16 `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-required:"true" env-default:"disable"`
}

type Migrations struct {
	Path string `yaml:"path" env-required:"true"`
}

func MustLoadConfig(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
