package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Postgres string `json:"postgres"`
	Redis    string `json:"redis"`
	RabbitMQ string `json:"rabbitmq"`
}

func HealthHandler(db *pgx.Conn, rdb *redis.Client, rabbitConn *amqp.Connection) http.HandlerFunc {
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

		rabbitStatus := "ok"
		if rabbitConn.IsClosed() {
			rabbitStatus = "down"
		}

		resp := HealthResponse{
			Status:   "ok",
			Postgres: pgStatus,
			Redis:    redisStatus,
			RabbitMQ: rabbitStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
