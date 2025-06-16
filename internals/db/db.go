package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func DBConnect() error {
	_ = godotenv.Load("../.env")
	url := os.Getenv("DB_URL")
	
	var err error
	Pool, err = pgxpool.New(context.Background(), url)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	return nil
}
