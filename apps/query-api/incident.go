package main

import "time"

type Incident struct {
	ID          string                 `json:"id"`
	ServiceName string                 `json:"service_name"`
	Type        string                 `json:"type"`
	Status      string                 `json:"status"`
	Title       string                 `json:"title"`
	Fingerprint *string                `json:"fingerprint,omitempty"`
	Severity    string                 `json:"severity"`
	StartedAt   time.Time              `json:"started_at"`
	LastSeenAt  time.Time              `json:"last_seen_at"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	LogCount    int                    `json:"log_count"`
	Context     map[string]interface{} `json:"context"`
}
