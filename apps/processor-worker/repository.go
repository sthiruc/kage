package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func InsertLog(db *pgx.Conn, logRecord Log) error {
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
		logRecord.ID,
		logRecord.Source,
		logRecord.ServiceName,
		logRecord.Environment,
		logRecord.Level,
		logRecord.Message,
		logRecord.Timestamp,
		logRecord.Metadata,
	)

	return err
}
