package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Postgres string `json:"postgres"`
	Redis    string `json:"redis"`
}

func HealthHandler(db *pgx.Conn, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		pgStatus := "ok"
		if err := db.Ping(ctx); err != nil {
			pgStatus = "down"
		}

		redisStatus := "ok"
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			redisStatus = "down"
		}

		resp := HealthResponse{
			Status:   "ok",
			Postgres: pgStatus,
			Redis:    redisStatus,
		}
		json.NewEncoder(w).Encode(resp)

	}
}
