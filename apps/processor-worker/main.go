package main

import (
	"context"
	"log"
)

func main() {
	cfg := LoadConfig()

	db := ConnectDB(cfg.PostgresURL)
	defer db.Close(context.Background())

	rabbitConn := ConnectRabbitMQ(cfg.RabbitMQURL)
	defer rabbitConn.Close()

	rabbitCh := OpenChannel(rabbitConn)
	defer rabbitCh.Close()

	DeclareQueue(rabbitCh, cfg.QueueName)

	log.Println("processor-worker started")

	if err := StartConsumer(rabbitCh, db, cfg.QueueName); err != nil {
		log.Fatalf("consumer failed: %v", err)
	}
}
