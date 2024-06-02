package env

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	// Local Configuration
	LocalGRPCPort         string `env:"LOCAL_GRPC_PORT"`
	LocalGRPCTimeout      string `env:"LOCAL_GRPC_TIMEOUT"`
	LocalPostgresHost     string `env:"LOCAL_POSTGRES_HOST"`
	LocalPostgresPort     string `env:"LOCAL_POSTGRES_PORT"`
	LocalPostgresUser     string `env:"LOCAL_POSTGRES_USER"`
	LocalPostgresPassword string `env:"LOCAL_POSTGRES_PASSWORD"`
	LocalPostgresDBName   string `env:"LOCAL_POSTGRES_DBNAME"`
	LocalPostgresSSLMode  string `env:"LOCAL_POSTGRES_SSLMODE"`
	LocalMigrationsPath   string `env:"LOCAL_MIGRATIONS_PATH"`

	// Production Configuration
	ProdGRPCPort         string `env:"PROD_GRPC_PORT"`
	ProdGRPCTimeout      string `env:"PROD_GRPC_TIMEOUT"`
	ProdPostgresHost     string `env:"PROD_POSTGRES_HOST"`
	ProdPostgresPort     string `env:"PROD_POSTGRES_PORT"`
	ProdPostgresUser     string `env:"PROD_POSTGRES_USER"`
	ProdPostgresPassword string `env:"PROD_POSTGRES_PASSWORD"`
	ProdPostgresDBName   string `env:"PROD_POSTGRES_DBNAME"`
	ProdPostgresSSLMode  string `env:"PROD_POSTGRES_SSLMODE"`
	ProdMigrationsPath   string `env:"PROD_MIGRATIONS_PATH"`
}

func MustLoadEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		panic(".env file does not exists")
	}

	var cfg Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
