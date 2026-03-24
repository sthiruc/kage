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

	rdb := ConnectRedis(cfg.RedisAddr)

	rabbitConn := ConnectRabbitMQ(cfg.RabbitMQURL)
	defer rabbitConn.Close()

	rabbitCh := OpenChannel(rabbitConn)
	defer rabbitCh.Close()

	DeclareQueue(rabbitCh, cfg.QueueName)

	r := chi.NewRouter()

	r.Get("/health", HealthHandler(db, rdb, rabbitConn))
	r.Post("/api/v1/logs", CreateLogHandler(rabbitCh, cfg.QueueName))

	log.Println("ingestion-api running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
