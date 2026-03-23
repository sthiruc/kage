package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func ConnectDB(url string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	return conn
}
