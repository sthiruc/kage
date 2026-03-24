package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func InsertLog(db *pgx.Conn, log Log) error {
	query := `
	INSERT INTO logs (
		id,
		source,
		service_name,
		environment,
		level,
		message,
		timestamp,
		metadata
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := db.Exec(context.Background(), query,
		log.ID,
		log.Source,
		log.ServiceName,
		log.Environment,
		log.Level,
		log.Message,
		log.Timestamp,
		log.Metadata,
	)

	return err
}
