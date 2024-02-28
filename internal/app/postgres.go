package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxConnDB = 500
)

func NewPoolPG(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsnPG())

	cfg.MaxConns = maxConnDB

	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new postgres pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping new pool: %w", err)
	}
	return pool, err
}

// instead of .env
func dsnPG() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		"usr", "pswd",
		"hst", 5432,
		"reqDB")
}
