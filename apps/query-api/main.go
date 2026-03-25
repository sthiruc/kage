package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := LoadConfig()

	db := ConnectDB(cfg.PostgresURL)
	defer db.Close(context.Background())

	r := chi.NewRouter()

	r.Get("/health", HealthHandler(db))
	r.Get("/api/v1/logs", GetLogsHandler(db))
	r.Get("/api/v1/logs/{id}", GetLogByIDHandler(db))
	r.Get("/api/v1/incidents", GetIncidentsHandler(db))
	r.Get("/api/v1/incidents/{id}", GetIncidentByIDHandler(db))
	r.Post("/api/v1/incidents/{id}/ack", AcknowledgeIncidentHandler(db))
	r.Post("/api/v1/incidents/{id}/resolve", ResolveIncidentHandler(db))

	addr := ":" + cfg.Port
	log.Printf("query-api running on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
