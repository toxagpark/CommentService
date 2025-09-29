package pkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPostgresConnect = errors.New("postgres connect error")
	ErrPostgresPing    = errors.New("postgres ping error")
)

func NewPostgres(url string, ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPostgresConnect, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPostgresPing, err)
	}

	return pool, nil
}

func ClosePostgres(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
