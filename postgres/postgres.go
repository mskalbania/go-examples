package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

var connection *pgx.Conn

func init() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("error pinging database=%v: %v", conn.Config().Database, err)
	}
	log.Printf("connected to database=%v", conn.Config().Database)
	connection = conn
}

func closeConnection() {
	err := connection.Close(context.Background())
	if err != nil {
		log.Fatalf("unable to close connection to postgres: %v", err)
	}
	log.Printf("closed connection to postgres")
}
