package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type sslMode string

const (
	SSLDisable    sslMode = "disable"
	SSLRequire    sslMode = "require"
	SSLVerifyFull sslMode = "verify-full"
	SSLVerifyCa   sslMode = "verify-ca"
)

func NewDB(host string, port int, user string, password string, dbname string, sslMode sslMode) (*gorm.DB, error) {
	const op = "database.connection.postgres.NewDB"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		host, user, password, port, dbname, sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
