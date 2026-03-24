package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateLogHandler(ch *amqp.Channel, queueName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logReq Log

		if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		logReq.ID = uuid.New().String()

		if err := PublishLog(ch, queueName, logReq); err != nil {
			http.Error(w, "failed to publish log", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "accepted",
			"log_id": logReq.ID,
		})
	}
}
