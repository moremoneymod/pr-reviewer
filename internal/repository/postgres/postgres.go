package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pgxPool *pgxpool.Pool
}

func New(ctx context.Context, pgConfig string) (*Storage, error) {

	pgCfg, err := pgxpool.ParseConfig(pgConfig)

	dbc, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, err
	}
	return &Storage{dbc}, nil
}
