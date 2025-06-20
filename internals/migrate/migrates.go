package migrate

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const migrationDir = "internals/migrations"

func RunMigrations(dsn string) error {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	db := stdlib.OpenDB(*config)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	log.Println("Running migrations...")
	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations applied successfully.")
	return nil
}
