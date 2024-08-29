package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go-examples/rest/config"
	"time"
)

type PostgresDatabase struct {
	Conn *pgx.Conn
}

func NewPostgresDatabase(config *config.AppConfig) (*PostgresDatabase, func(), error) {
	conn, err := pgx.Connect(context.TODO(), connectionString(config))
	if err != nil {
		return nil, nil, err
	}
	if err := conn.Ping(context.TODO()); err != nil {
		return nil, nil, err
	}
	return &PostgresDatabase{Conn: conn}, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn.Close(ctx)
	}, nil
}

func connectionString(config *config.AppConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Database)
}
