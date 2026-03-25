package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Incident struct {
	ID          string
	ServiceName string
	Type        string
	Status      string
	Title       string
	Fingerprint *string
	Severity    string
	StartedAt   time.Time
	LastSeenAt  time.Time
	ResolvedAt  *time.Time
	LogCount    int
	Context     map[string]interface{}
}

func GetOpenIncidentByServiceAndType(db *pgx.Conn, serviceName, incidentType string) (*Incident, error) {
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
		WHERE service_name = $1
		  AND type = $2
		  AND status = 'open'
		ORDER BY started_at DESC
		LIMIT 1
	`

	var incident Incident

	err := db.QueryRow(context.Background(), query, serviceName, incidentType).Scan(
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

func CreateIncident(db *pgx.Conn, serviceName, incidentType, title, severity string, contextData map[string]interface{}) error {
	query := `
		INSERT INTO incidents (
			id,
			service_name,
			type,
			status,
			title,
			severity,
			started_at,
			last_seen_at,
			log_count,
			context
		)
		VALUES ($1, $2, $3, 'open', $4, $5, $6, $7, $8, $9)
	`

	now := time.Now().UTC()

	_, err := db.Exec(
		context.Background(),
		query,
		uuid.New().String(),
		serviceName,
		incidentType,
		title,
		severity,
		now,
		now,
		1,
		contextData,
	)

	return err
}

func UpdateIncidentLastSeen(db *pgx.Conn, incidentID string) error {
	query := `
		UPDATE incidents
		SET
			last_seen_at = $2,
			log_count = log_count + 1
		WHERE id = $1
	`

	_, err := db.Exec(context.Background(), query, incidentID, time.Now().UTC())
	return err
}
