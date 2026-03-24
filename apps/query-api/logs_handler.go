package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func GetLogsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		offset := 0

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			parsed, err := strconv.Atoi(limitStr)
			if err != nil || parsed <= 0 {
				http.Error(w, "invalid limit", http.StatusBadRequest)
				return
			}
			if parsed > 200 {
				parsed = 200
			}
			limit = parsed
		}

		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			parsed, err := strconv.Atoi(offsetStr)
			if err != nil || parsed < 0 {
				http.Error(w, "invalid offset", http.StatusBadRequest)
				return
			}
			offset = parsed
		}

		filters := LogFilters{
			ServiceName: r.URL.Query().Get("service_name"),
			Level:       r.URL.Query().Get("level"),
			Limit:       limit,
			Offset:      offset,
		}

		logs, err := GetLogs(db, filters)
		if err != nil {
			http.Error(w, "failed to fetch logs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

func GetLogByIDHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		logRecord, err := GetLogByID(db, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "log not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to fetch log", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logRecord)
	}
}
