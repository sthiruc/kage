package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

const ErrorSpikeIncidentType = "error_spike"

func HandleIncidentDetection(
	db *pgx.Conn,
	rdb *redis.Client,
	logRecord Log,
	threshold int,
	windowSec int,
) error {
	if logRecord.Level != "ERROR" {
		return nil
	}

	now := time.Now().UTC()
	bucket := now.Unix() / int64(windowSec)

	key := fmt.Sprintf("incident:error_spike:%s:%d", logRecord.ServiceName, bucket)

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	err = rdb.Expire(ctx, key, time.Duration(windowSec*2)*time.Second).Err()
	if err != nil {
		return err
	}

	log.Printf("error counter for service=%s bucket=%d count=%d", logRecord.ServiceName, bucket, count)

	if int(count) < threshold {
		return nil
	}

	existing, err := GetOpenIncidentByServiceAndType(db, logRecord.ServiceName, ErrorSpikeIncidentType)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	if existing != nil {
		return UpdateIncidentLastSeen(db, existing.ID)
	}

	contextData := map[string]interface{}{
		"threshold":        threshold,
		"window_seconds":   windowSec,
		"service_name":     logRecord.ServiceName,
		"triggered_by_log": logRecord.ID,
		"count":            count,
	}

	title := fmt.Sprintf("High ERROR rate detected in service %s", logRecord.ServiceName)

	return CreateIncident(
		db,
		logRecord.ServiceName,
		ErrorSpikeIncidentType,
		title,
		"high",
		contextData,
	)
}
