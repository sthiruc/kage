package main

import (
	"os"
	"strconv"
)

type Config struct {
	PostgresURL         string
	RabbitMQURL         string
	RedisAddr           string
	QueueName           string
	ErrorSpikeThreshold int
	ErrorSpikeWindowSec int
}

func LoadConfig() Config {
	return Config{
		PostgresURL:         getEnv("POSTGRES_URL", "postgres://logging_user:logging_pass@localhost:5432/logging"),
		RabbitMQURL:         getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		QueueName:           getEnv("QUEUE_NAME", "logs_ingestion"),
		ErrorSpikeThreshold: getEnvInt("ERROR_SPIKE_THRESHOLD", 5),
		ErrorSpikeWindowSec: getEnvInt("ERROR_SPIKE_WINDOW_SEC", 300),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return parsed
}
