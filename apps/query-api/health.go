package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Postgres string `json:"postgres"`
}

func HealthHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pgStatus := "ok"
		if err := db.Ping(context.Background()); err != nil {
			pgStatus = "down"
		}

		resp := HealthResponse{
			Status:   "ok",
			Postgres: pgStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
