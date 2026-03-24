package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type LogFilters struct {
	ServiceName string
	Level       string
	Limit       int
	Offset      int
}

func GetLogs(db *pgx.Conn, filters LogFilters) ([]Log, error) {
	baseQuery := `
		SELECT
			id,
			source,
			service_name,
			environment,
			level,
			message,
			fingerprint,
			trace_id,
			request_id,
			host,
			timestamp,
			received_at,
			metadata
		FROM logs
	`

	var conditions []string
	var args []interface{}
	argPos := 1

	if filters.ServiceName != "" {
		conditions = append(conditions, "service_name = $"+strconv.Itoa(argPos))
		args = append(args, filters.ServiceName)
		argPos++
	}

	if filters.Level != "" {
		conditions = append(conditions, "level = $"+strconv.Itoa(argPos))
		args = append(args, filters.Level)
		argPos++
	}

	query := baseQuery

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY timestamp DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []Log

	for rows.Next() {
		var logRecord Log

		err := rows.Scan(
			&logRecord.ID,
			&logRecord.Source,
			&logRecord.ServiceName,
			&logRecord.Environment,
			&logRecord.Level,
			&logRecord.Message,
			&logRecord.Fingerprint,
			&logRecord.TraceID,
			&logRecord.RequestID,
			&logRecord.Host,
			&logRecord.Timestamp,
			&logRecord.ReceivedAt,
			&logRecord.Metadata,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, logRecord)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func GetLogByID(db *pgx.Conn, id string) (*Log, error) {
	query := `
		SELECT
			id,
			source,
			service_name,
			environment,
			level,
			message,
			fingerprint,
			trace_id,
			request_id,
			host,
			timestamp,
			received_at,
			metadata
		FROM logs
		WHERE id = $1
	`

	var logRecord Log

	err := db.QueryRow(context.Background(), query, id).Scan(
		&logRecord.ID,
		&logRecord.Source,
		&logRecord.ServiceName,
		&logRecord.Environment,
		&logRecord.Level,
		&logRecord.Message,
		&logRecord.Fingerprint,
		&logRecord.TraceID,
		&logRecord.RequestID,
		&logRecord.Host,
		&logRecord.Timestamp,
		&logRecord.ReceivedAt,
		&logRecord.Metadata,
	)
	if err != nil {
		return nil, err
	}

	return &logRecord, nil
}
