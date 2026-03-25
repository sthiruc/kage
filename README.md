Kage — Distributed Logging & Incident Detection System

Kage is a backend system that ingests logs, processes them asynchronously, stores them, and detects incidents in real time.

⸻

How It Works
	1.	A client sends a log to the ingestion API
	2.	The ingestion API validates and publishes the log to RabbitMQ
	3.	The processor worker consumes the message
	4.	The worker:
	•	stores the log in Postgres
	•	updates Redis counters
	•	runs incident detection logic
	5.	If a rule is triggered → an incident is created/updated
	6.	The query API allows fetching logs and incidents

Client → ingestion-api → RabbitMQ → processor-worker → Postgres
Client → query-api → Postgres
Worker → Redis (ephemeral detection state)


⸻

Architecture

Services
	•	ingestion-api
	•	Accepts logs via HTTP
	•	Publishes messages to RabbitMQ
	•	processor-worker
	•	Consumes messages
	•	Stores logs in Postgres
	•	Runs incident detection
	•	query-api
	•	Fetch logs with filters
	•	Fetch incidents
	•	Manage incident lifecycle (ack/resolve)

Infrastructure
	•	Postgres → durable storage
	•	Redis → counters + short-lived state
	•	RabbitMQ → async queue

⸻

Features

Log Ingestion
	•	Accept structured logs via HTTP
	•	Async processing (non-blocking ingestion)

Log Querying
	•	Filter by:
	•	service_name
	•	level
	•	time range
	•	Pagination (limit/offset)

Incident Detection
	•	Detect error spikes per service
	•	Uses Redis counters with time buckets

Incident Management
	•	List incidents
	•	Acknowledge incidents
	•	Resolve incidents

⸻

Getting Started

Prerequisites
	•	Docker
	•	Docker Compose

Run everything

docker compose up --build

Services
	•	ingestion-api → http://localhost:8080
	•	query-api → http://localhost:8081
	•	RabbitMQ UI → http://localhost:15672 (guest/guest)

⸻

API Reference

⸻

Ingestion API (8080)

Create Log

POST /api/v1/logs

Example Request

curl -X POST http://localhost:8080/api/v1/logs \
  -H "Content-Type: application/json" \
  -d '{
    "source": "payment-service",
    "service_name": "payments",
    "environment": "dev",
    "level": "ERROR",
    "message": "database timeout",
    "timestamp": "2026-03-25T12:00:00Z",
    "metadata": {
      "user_id": "123",
      "route": "/charge"
    }
  }'

Response

{
  "status": "accepted",
  "log_id": "uuid"
}


⸻

Query API (8081)

Get Logs

GET /api/v1/logs

Example Queries

# All logs
curl http://localhost:8081/api/v1/logs

# Filter by service
curl "http://localhost:8081/api/v1/logs?service_name=payments"

# Filter by level
curl "http://localhost:8081/api/v1/logs?level=ERROR"

# Time range
curl "http://localhost:8081/api/v1/logs?start=2026-03-25T12:00:00Z&end=2026-03-25T13:00:00Z"

# Combined
curl "http://localhost:8081/api/v1/logs?service_name=payments&level=ERROR&limit=10"


⸻

Get Log by ID

GET /api/v1/logs/{id}

curl http://localhost:8081/api/v1/logs/<LOG_ID>


⸻

Incidents

Get Incidents

GET /api/v1/incidents

# All incidents
curl http://localhost:8081/api/v1/incidents

# Only open
curl "http://localhost:8081/api/v1/incidents?status=open"

# By service
curl "http://localhost:8081/api/v1/incidents?service_name=payments"


⸻

Get Incident by ID

GET /api/v1/incidents/{id}

curl http://localhost:8081/api/v1/incidents/<INCIDENT_ID>


⸻

Acknowledge Incident

POST /api/v1/incidents/{id}/ack

curl -X POST http://localhost:8081/api/v1/incidents/<ID>/ack


⸻

Resolve Incident

POST /api/v1/incidents/{id}/resolve

curl -X POST http://localhost:8081/api/v1/incidents/<ID>/resolve


⸻

🏁 Summary

Kage demonstrates:
	•	async processing with queues
	•	distributed system design
	•	incident detection logic
	•	service separation and scalability


⸻

License

MIT