package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/xerrors"
)

func WithConfig(ctx context.Context, url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, xerrors.Errorf("Couldn't parse pgx config string: %w", err)
	}

	config.ConnConfig.LogLevel = pgx.LogLevelTrace
	config.MaxConns = 20

	return pgxpool.ConnectConfig(ctx, config)
}
