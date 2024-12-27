package db

import (
	"context"
	"fmt"

	"backend-layout/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPgx(conf *config.Config, ctx context.Context) (*pgxpool.Pool, error) {

	pgxConfig, err := pgxpool.ParseConfig(conf.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pgxConfig.MinConns = int32(conf.DB.MaxIdleConns)
	pgxConfig.MaxConns = int32(conf.DB.MaxOpenConns)

	dbpool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return dbpool, nil
}
