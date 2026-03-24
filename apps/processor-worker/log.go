package main

import "time"

type Log struct {
	ID          string                 `json:"id"`
	Source      string                 `json:"source"`
	ServiceName string                 `json:"service_name"`
	Environment string                 `json:"environment"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}
