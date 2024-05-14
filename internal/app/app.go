package app

import (
	"fmt"
	grpcapp "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/app/grpc"
	"github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/config"
	pgconnection "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/database/connection/postgres"
	pgmigration "github.com/g-vinokurov/pyramidum-backend-service-tasks/internal/database/migration/postgres"
	"log/slog"
	"os"
)

type App struct {
	grpcApp  *grpcapp.App
	migrator *pgmigration.Migrator
	database *pgconnection.Database
}

func NewApp(log *slog.Logger, cfg *config.Config) (*App, error) {
	const op = "app.NewApp"

	database, err := pgconnection.NewDatabase(
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.DBName,
		cfg.PostgreSQL.SSLMode,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("database connected")

	migrationFiles := os.DirFS(cfg.Migrations.Path)

	migrator, err := pgmigration.NewMigrator(migrationFiles, ".")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("migrator created")

	err = migrator.ApplyMigrations(database.DB(), cfg.PostgreSQL.DBName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("migrations applied")

	grpcApp := grpcapp.NewApp(log, int(cfg.GRPC.Port), database)

	return &App{
		grpcApp:  grpcApp,
		migrator: migrator,
		database: database,
	}, nil
}

func (a *App) Run() error {
	return a.grpcApp.Run()
}

func (a *App) Stop() {
	defer func() {
		_ = a.migrator.Close()
		_ = a.database.DB().Close()
	}()
	a.grpcApp.Stop()
}
