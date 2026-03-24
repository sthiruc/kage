package main

import (
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer(ch *amqp.Channel, db *pgx.Conn, queueName string) error {
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Printf("waiting for messages from queue %s", queueName)

	for msg := range msgs {
		var logRecord Log

		if err := json.Unmarshal(msg.Body, &logRecord); err != nil {
			log.Printf("failed to decode message: %v", err)
			_ = msg.Nack(false, false)
			continue
		}

		if err := InsertLog(db, logRecord); err != nil {
			log.Printf("failed to insert log %s: %v", logRecord.ID, err)
			_ = msg.Nack(false, true)
			continue
		}

		log.Printf("processed log %s", logRecord.ID)
		_ = msg.Ack(false)
	}

	return nil
}
