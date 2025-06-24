package migrate

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	migrationDir = "internals/migrations"
	maxRetries   = 10
	retryDelay   = 2 * time.Second
)

func RunMigrationsWithRetry(dsn string) error {
	var err error
	for i := 1; i <= maxRetries; i++ {
		err = run(dsn)
		if err == nil {
			return nil
		}

		log.Printf("Migration attempt %d failed: %v", i, err)
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("failed to run migrations after %d attempts: %w", maxRetries, err)
}

func run(dsn string) error {
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