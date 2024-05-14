package postgres

import (
	"database/sql"
	"fmt"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(host string, port uint16, user string, password string, dbname string, sslMode string) (*Database, error) {
	const op = "database.connection.postgres.NewDB"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		host, user, password, dbname, port, sslMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Database{db: db}, nil
}

func (d *Database) DB() *sql.DB {
	return d.db
}
