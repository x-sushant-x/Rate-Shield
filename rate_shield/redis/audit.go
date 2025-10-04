package redisClient

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
)

const (
	auditLogsKey = "audit:logs"
)

// RedisAudit implements the RedisAuditClient interface
type RedisAudit struct {
	client *redis.Client
}

// NewAuditClient creates a new Redis audit client using the existing rules client connection
func NewAuditClient(client *redis.Client) RedisAuditClient {
	return RedisAudit{
		client: client,
	}
}

// AppendAuditLog adds a new audit log entry to the Redis list (append-only)
func (r RedisAudit) AppendAuditLog(auditLog models.AuditLog) error {
	// Marshal audit log to JSON
	logJSON, err := json.Marshal(auditLog)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal audit log to JSON")
		return err
	}

	// Append to Redis list using RPUSH (right push for append-only behavior)
	err = r.client.RPush(context.Background(), auditLogsKey, string(logJSON)).Err()
	if err != nil {
		log.Error().Err(err).Msg("failed to append audit log to Redis")
		return err
	}

	log.Debug().
		Str("action", auditLog.Action).
		Str("actor", auditLog.Actor).
		Str("endpoint", auditLog.Endpoint).
		Msg("audit log appended successfully")

	return nil
}

// GetAuditLogs retrieves audit logs from Redis list within the specified range
// start and end are 0-based indices (0 is the first element, -1 is the last)
func (r RedisAudit) GetAuditLogs(start, end int64) ([]models.AuditLog, error) {
	// Get logs from Redis list using LRANGE
	results, err := r.client.LRange(context.Background(), auditLogsKey, start, end).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve audit logs from Redis")
		return nil, err
	}

	// Parse JSON strings to AuditLog structs
	logs := make([]models.AuditLog, 0, len(results))
	for _, result := range results {
		var auditLog models.AuditLog
		err := json.Unmarshal([]byte(result), &auditLog)
		if err != nil {
			log.Error().Err(err).Str("log_data", result).Msg("failed to unmarshal audit log")
			continue // Skip malformed entries
		}
		logs = append(logs, auditLog)
	}

	return logs, nil
}

// GetAuditLogCount returns the total number of audit logs stored
func (r RedisAudit) GetAuditLogCount() (int64, error) {
	count, err := r.client.LLen(context.Background(), auditLogsKey).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to get audit log count from Redis")
		return 0, err
	}

	return count, nil
}

// GetAllAuditLogs retrieves all audit logs (wrapper for convenience)
func (r RedisAudit) GetAllAuditLogs() ([]models.AuditLog, error) {
	return r.GetAuditLogs(0, -1)
}
