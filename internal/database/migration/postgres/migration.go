package postgres

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"
	"io/fs"
)

type Migrator struct {
	srcDriver source.Driver
}

func NewMigrator(sqlFiles fs.FS, dirName string) (*Migrator, error) {
	const op = "database.migration.postgres.NewMigrator"

	driver, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Migrator{srcDriver: driver}, nil
}

func (m *Migrator) ApplyMigrations(db *gorm.DB, dbName string) error {
	const op = "database.migration.postgres.ApplyMigrations"

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	migrator, err := migrate.NewWithInstance("migration_embedded_sql_files", m.srcDriver, dbName, driver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_, _ = migrator.Close()
	}()

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: unable to apply migration: %w", op, err)
	}

	return nil
}
