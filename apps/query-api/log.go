package main

import "time"

type Log struct {
	ID          string                 `json:"id"`
	Source      string                 `json:"source"`
	ServiceName string                 `json:"service_name"`
	Environment string                 `json:"environment"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Fingerprint *string                `json:"fingerprint,omitempty"`
	TraceID     *string                `json:"trace_id,omitempty"`
	RequestID   *string                `json:"request_id,omitempty"`
	Host        *string                `json:"host,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	ReceivedAt  time.Time              `json:"received_at"`
	Metadata    map[string]interface{} `json:"metadata"`
}
