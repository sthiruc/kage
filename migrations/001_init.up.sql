CREATE TABLE IF NOT EXISTS logs (
    id UUID PRIMARY KEY,
    source VARCHAR(100) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    environment VARCHAR(50) NOT NULL,
    level VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
    fingerprint VARCHAR(255),
    trace_id VARCHAR(100),
    request_id VARCHAR(100),
    host VARCHAR(255),
    timestamp TIMESTAMPTZ NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_logs_service_timestamp
    ON logs(service_name, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_logs_level_timestamp
    ON logs(level, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_logs_request_id
    ON logs(request_id);

CREATE INDEX IF NOT EXISTS idx_logs_trace_id
    ON logs(trace_id);

CREATE INDEX IF NOT EXISTS idx_logs_fingerprint_timestamp
    ON logs(fingerprint, timestamp DESC);

CREATE TABLE IF NOT EXISTS incidents (
    id UUID PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(30) NOT NULL,
    title VARCHAR(255) NOT NULL,
    fingerprint VARCHAR(255),
    severity VARCHAR(30) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ NOT NULL,
    resolved_at TIMESTAMPTZ,
    log_count INT NOT NULL DEFAULT 0,
    context JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_incidents_service_status
    ON incidents(service_name, status);

CREATE INDEX IF NOT EXISTS idx_incidents_status_last_seen
    ON incidents(status, last_seen_at DESC);

CREATE INDEX IF NOT EXISTS idx_incidents_fingerprint_status
    ON incidents(fingerprint, status);