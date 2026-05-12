package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"graft/server/internal/config"
)

// Open creates the PostgreSQL GORM connection required by the server runtime.
func Open(cfg config.DatabaseConfig) (*gorm.DB, error) {
	if cfg.Driver != "postgres" {
		return nil, fmt.Errorf("unsupported database driver %q: only postgres is supported", cfg.Driver)
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres database: %w", err)
	}

	return db, nil
}
