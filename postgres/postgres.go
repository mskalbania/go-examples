package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

var connection *pgx.Conn

func openConnection(user string, password string, host string, port int, database string) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, database)
	conn, err := pgx.Connect(context.Background(), url)
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
