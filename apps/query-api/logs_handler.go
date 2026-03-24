package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type LogsListResponse struct {
	Logs   []Log `json:"logs"`
	Count  int   `json:"count"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

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

		var startPtr *time.Time
		if startStr := r.URL.Query().Get("start"); startStr != "" {
			parsed, err := time.Parse(time.RFC3339, startStr)
			if err != nil {
				http.Error(w, "invalid start, must be RFC3339", http.StatusBadRequest)
				return
			}
			startPtr = &parsed
		}

		var endPtr *time.Time
		if endStr := r.URL.Query().Get("end"); endStr != "" {
			parsed, err := time.Parse(time.RFC3339, endStr)
			if err != nil {
				http.Error(w, "invalid end, must be RFC3339", http.StatusBadRequest)
				return
			}
			endPtr = &parsed
		}

		if startPtr != nil && endPtr != nil && startPtr.After(*endPtr) {
			http.Error(w, "start must be before or equal to end", http.StatusBadRequest)
			return
		}

		filters := LogFilters{
			ServiceName: r.URL.Query().Get("service_name"),
			Level:       r.URL.Query().Get("level"),
			Start:       startPtr,
			End:         endPtr,
			Limit:       limit,
			Offset:      offset,
		}

		logs, err := GetLogs(db, filters)
		if err != nil {
			http.Error(w, "failed to fetch logs", http.StatusInternalServerError)
			return
		}

		response := LogsListResponse{
			Logs:   logs,
			Count:  len(logs),
			Limit:  limit,
			Offset: offset,
		}

		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(response)

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
