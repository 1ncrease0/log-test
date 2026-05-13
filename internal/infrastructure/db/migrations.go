package db

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func (db *DB) Migrate() error {
	db.log.Debug("running migration")
	files, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("migration source: %w", err)
	}
	driver, err := pgxmigrate.WithInstance(db.conn.DB, &pgxmigrate.Config{})
	if err != nil {
		return fmt.Errorf("migration pgx driver: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", files, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migration instance: %w", err)
	}

	err = m.Up()

	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			db.log.Error("migration", "error", err)
			return fmt.Errorf("migration up: %w", err)
		}
		db.log.Debug("migration did not change anything")
	}

	db.log.Debug("migration finished")
	return nil
}
