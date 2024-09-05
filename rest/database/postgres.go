package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-examples/rest/config"
	"time"
)

type Database interface {
	Ping(context.Context) error
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewPostgresDatabase(config *config.AppConfig) (Database, func(), error) {
	pool, err := pgxpool.New(context.Background(), connectionString(config))
	if err != nil {
		return nil, nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, nil, err
	}
	return pool, func() {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		pool.Close()
	}, nil
}

func connectionString(config *config.AppConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&pool_min_conns=%d",
		config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Database, config.DB.PoolMax, config.DB.PoolMin)
}
