package main

import "os"

type Config struct {
	PostgresURL string
	Port        string
}

func LoadConfig() Config {
	return Config{
		PostgresURL: getEnv("POSTGRES_URL", "postgres://logging_user:logging_pass@localhost:5432/logging"),
		Port:        getEnv("PORT", "8081"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
