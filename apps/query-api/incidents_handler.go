package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func GetIncidentsHandler(db *pgx.Conn) http.HandlerFunc {
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

		filters := IncidentFilters{
			ServiceName: r.URL.Query().Get("service_name"),
			Status:      r.URL.Query().Get("status"),
			Type:        r.URL.Query().Get("type"),
			Limit:       limit,
			Offset:      offset,
		}

		incidents, err := GetIncidents(db, filters)
		if err != nil {
			http.Error(w, "failed to fetch incidents", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(incidents)
	}
}

func GetIncidentByIDHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		incident, err := GetIncidentByID(db, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "incident not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to fetch incident", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(incident)
	}
}

func AcknowledgeIncidentHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		err := AcknowledgeIncident(db, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "incident not found or not open", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to acknowledge incident", http.StatusInternalServerError)
			return
		}

		incident, err := GetIncidentByID(db, id)
		if err != nil {
			http.Error(w, "incident acknowledged but failed to fetch updated record", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(incident)
	}
}

func ResolveIncidentHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		err := ResolveIncident(db, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "incident not found or already resolved", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to resolve incident", http.StatusInternalServerError)
			return
		}

		incident, err := GetIncidentByID(db, id)
		if err != nil {
			http.Error(w, "incident resolved but failed to fetch updated record", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(incident)
	}
}
