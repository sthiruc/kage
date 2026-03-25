package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type IncidentFilters struct {
	ServiceName string
	Status      string
	Type        string
	Limit       int
	Offset      int
}

func GetIncidents(db *pgx.Conn, filters IncidentFilters) ([]Incident, error) {
	baseQuery := `
		SELECT
			id,
			service_name,
			type,
			status,
			title,
			fingerprint,
			severity,
			started_at,
			last_seen_at,
			resolved_at,
			log_count,
			context
		FROM incidents
	`

	var conditions []string
	var args []interface{}
	argPos := 1

	if filters.ServiceName != "" {
		conditions = append(conditions, "service_name = $"+strconv.Itoa(argPos))
		args = append(args, filters.ServiceName)
		argPos++
	}

	if filters.Status != "" {
		conditions = append(conditions, "status = $"+strconv.Itoa(argPos))
		args = append(args, filters.Status)
		argPos++
	}

	if filters.Type != "" {
		conditions = append(conditions, "type = $"+strconv.Itoa(argPos))
		args = append(args, filters.Type)
		argPos++
	}

	query := baseQuery

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY last_seen_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []Incident

	for rows.Next() {
		var incident Incident

		err := rows.Scan(
			&incident.ID,
			&incident.ServiceName,
			&incident.Type,
			&incident.Status,
			&incident.Title,
			&incident.Fingerprint,
			&incident.Severity,
			&incident.StartedAt,
			&incident.LastSeenAt,
			&incident.ResolvedAt,
			&incident.LogCount,
			&incident.Context,
		)
		if err != nil {
			return nil, err
		}

		incidents = append(incidents, incident)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}

func GetIncidentByID(db *pgx.Conn, id string) (*Incident, error) {
	query := `
		SELECT
			id,
			service_name,
			type,
			status,
			title,
			fingerprint,
			severity,
			started_at,
			last_seen_at,
			resolved_at,
			log_count,
			context
		FROM incidents
		WHERE id = $1
	`

	var incident Incident

	err := db.QueryRow(context.Background(), query, id).Scan(
		&incident.ID,
		&incident.ServiceName,
		&incident.Type,
		&incident.Status,
		&incident.Title,
		&incident.Fingerprint,
		&incident.Severity,
		&incident.StartedAt,
		&incident.LastSeenAt,
		&incident.ResolvedAt,
		&incident.LogCount,
		&incident.Context,
	)
	if err != nil {
		return nil, err
	}

	return &incident, nil
}

func AcknowledgeIncident(db *pgx.Conn, id string) error {
	query := `
		UPDATE incidents
		SET status = 'acknowledged'
		WHERE id = $1
		  AND status = 'open'
	`

	commandTag, err := db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func ResolveIncident(db *pgx.Conn, id string) error {
	query := `
		UPDATE incidents
		SET
			status = 'resolved',
			resolved_at = $2
		WHERE id = $1
		  AND status IN ('open', 'acknowledged')
	`

	commandTag, err := db.Exec(context.Background(), query, id, time.Now().UTC())
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
