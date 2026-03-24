package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreateLogHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logReq Log

		if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		logReq.ID = uuid.New().String()

		if err := InsertLog(db, logReq); err != nil {
			http.Error(w, "failed to insert log", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(logReq)
	}
}
