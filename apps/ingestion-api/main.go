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

	r := chi.NewRouter()

	r.Get("/health", HealthHandler(db, rdb))

	log.Println("ingestion-api running on :8080")
	http.ListenAndServe(":8080", r)
}
