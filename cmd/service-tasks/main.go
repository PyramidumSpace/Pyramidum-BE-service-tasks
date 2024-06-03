package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/app"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/config"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/env"
	slogattr "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/lib/log/slog/attr"
)

func main() {
	log := slog.Default()

	// init environment variables from .env file
	env.MustLoadEnv()

	log.Info("environment variables loaded")

	cfg := config.MustLoadConfig()

	log.Info("config loaded", slog.Any("config", cfg))

	a, err := app.NewApp(log, cfg)
	if err != nil {
		log.Error("cannot create app", slogattr.Err(err))
		os.Exit(1)
	}

	go func() {
		if err := a.Run(); err != nil {
			log.Error("cannot run app", slogattr.Err(err))
			os.Exit(1)
		}
	}()

	log.Info("app started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("stopping app")

	a.Stop()

	log.Info("app stopped")
}
