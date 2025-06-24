package db

import (
	"context"
	"fmt"


	"github.com/jackc/pgx/v5/pgxpool"

)

var Pool *pgxpool.Pool

func DBConnect(dsn string) error {
	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("pgxpool connect failed: %w", err)
	}
	return nil

}
