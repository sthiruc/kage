package main

import "os"

type Config struct {
	PostgresURL string
	RabbitMQURL string
	QueueName   string
}

func LoadConfig() Config {
	return Config{
		PostgresURL: getEnv("POSTGRES_URL", "postgres://logging_user:logging_pass@localhost:5432/logging"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName:   getEnv("QUEUE_NAME", "logs_ingestion"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
