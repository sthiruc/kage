package main

import (
	"context"
	"log"
)

func main() {
	cfg := LoadConfig()

	db := ConnectDB(cfg.PostgresURL)
	defer db.Close(context.Background())

	rdb := ConnectRedis(cfg.RedisAddr)
	defer rdb.Close()

	rabbitConn := ConnectRabbitMQ(cfg.RabbitMQURL)
	defer rabbitConn.Close()

	rabbitCh := OpenChannel(rabbitConn)
	defer rabbitCh.Close()

	DeclareQueue(rabbitCh, cfg.QueueName)

	log.Println("processor-worker started")

	if err := StartConsumer(
		rabbitCh,
		db,
		rdb,
		cfg.QueueName,
		cfg.ErrorSpikeThreshold,
		cfg.ErrorSpikeWindowSec,
	); err != nil {
		log.Fatalf("consumer failed: %v", err)
	}
}
