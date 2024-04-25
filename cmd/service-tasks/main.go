package main

import (
	"log/slog"
	"os"
	"pyramidum-backend-service-tasks/internal/config"
	"pyramidum-backend-service-tasks/internal/env"
)

func main() {
	// init environment variables from .env file
	env.MustLoadEnv()

	configPath := os.Getenv("CONFIG_PATH")
	cfg := config.MustLoadConfig(configPath)

	slog.Info("config loaded", slog.Any("config", cfg))
}
