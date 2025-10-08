package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
)

// AuditService defines the interface for audit logging operations
type AuditService interface {
	LogRuleChange(actor, action, endpoint string, oldRule, newRule *models.Rule, ipAddress, userAgent string) error
	GetAuditLogs(page, items int) (models.PaginatedAuditLogs, error)
	GetAllAuditLogs() ([]models.AuditLog, error)
	GetAuditLogsByEndpoint(endpoint string) ([]models.AuditLog, error)
	GetAuditLogsByActor(actor string) ([]models.AuditLog, error)
	GetAuditLogsByAction(action string) ([]models.AuditLog, error)
}

// AuditServiceRedis implements the AuditService interface using Redis
type AuditServiceRedis struct {
	auditClient redisClient.RedisAuditClient
}

// NewAuditService creates a new audit service instance
func NewAuditService(auditClient redisClient.RedisAuditClient) AuditService {
	return &AuditServiceRedis{
		auditClient: auditClient,
	}
}

// LogRuleChange logs a rule modification event to the audit trail
func (s *AuditServiceRedis) LogRuleChange(
	actor, action, endpoint string,
	oldRule, newRule *models.Rule,
	ipAddress, userAgent string,
) error {
	// Validate action
	if action != models.AuditActionCreate &&
		action != models.AuditActionUpdate &&
		action != models.AuditActionDelete {
		return errors.New("invalid audit action")
	}

	// Create audit log entry
	auditLog := models.AuditLog{
		ID:        uuid.New().String(),
		Timestamp: time.Now().Unix(),
		Actor:     actor,
		Action:    action,
		Endpoint:  endpoint,
		OldRule:   oldRule,
		NewRule:   newRule,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	// Save to Redis
	err := s.auditClient.AppendAuditLog(auditLog)
	if err != nil {
		log.Error().Err(err).
			Str("action", action).
			Str("actor", actor).
			Str("endpoint", endpoint).
			Msg("failed to log audit event")
		return err
	}

	log.Info().
		Str("id", auditLog.ID).
		Str("action", action).
		Str("actor", actor).
		Str("endpoint", endpoint).
		Msg("audit event logged successfully")

	return nil
}

// GetAllAuditLogs retrieves all audit logs from the system
func (s *AuditServiceRedis) GetAllAuditLogs() ([]models.AuditLog, error) {
	logs, err := s.auditClient.GetAllAuditLogs()
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve all audit logs")
		return nil, err
	}

	return logs, nil
}

// GetAuditLogs retrieves audit logs with pagination
func (s *AuditServiceRedis) GetAuditLogs(page, items int) (models.PaginatedAuditLogs, error) {
	// Get total count
	totalCount, err := s.auditClient.GetAuditLogCount()
	if err != nil {
		log.Error().Err(err).Msg("failed to get audit log count")
		return models.PaginatedAuditLogs{}, err
	}

	if totalCount == 0 {
		return models.PaginatedAuditLogs{
			PageNumber:  1,
			TotalItems:  0,
			HasNextPage: false,
			Logs:        []models.AuditLog{},
		}, nil
	}

	// Calculate start and end indices for pagination
	// Redis LRANGE uses 0-based indexing, and we want newest logs first
	// So we reverse the pagination: newest logs are at the end of the list
	start := totalCount - int64(page*items)
	end := start + int64(items) - 1

	if start < 0 {
		start = 0
	}
	if end >= totalCount {
		end = totalCount - 1
	}

	// Check if page is out of range
	if start > end {
		return models.PaginatedAuditLogs{}, errors.New("page number out of range")
	}

	// Get logs from Redis
	logs, err := s.auditClient.GetAuditLogs(start, end)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve paginated audit logs")
		return models.PaginatedAuditLogs{}, err
	}

	// Reverse logs to show newest first
	reversedLogs := make([]models.AuditLog, len(logs))
	for i, log := range logs {
		reversedLogs[len(logs)-1-i] = log
	}

	hasNextPage := start > 0

	return models.PaginatedAuditLogs{
		PageNumber:  page,
		TotalItems:  len(reversedLogs),
		HasNextPage: hasNextPage,
		Logs:        reversedLogs,
	}, nil
}

// GetAuditLogsByEndpoint filters audit logs by endpoint
func (s *AuditServiceRedis) GetAuditLogsByEndpoint(endpoint string) ([]models.AuditLog, error) {
	allLogs, err := s.GetAllAuditLogs()
	if err != nil {
		return nil, err
	}

	filteredLogs := make([]models.AuditLog, 0)
	for _, log := range allLogs {
		if log.Endpoint == endpoint {
			filteredLogs = append(filteredLogs, log)
		}
	}

	return filteredLogs, nil
}

// GetAuditLogsByActor filters audit logs by actor
func (s *AuditServiceRedis) GetAuditLogsByActor(actor string) ([]models.AuditLog, error) {
	allLogs, err := s.GetAllAuditLogs()
	if err != nil {
		return nil, err
	}

	filteredLogs := make([]models.AuditLog, 0)
	for _, log := range allLogs {
		if log.Actor == actor {
			filteredLogs = append(filteredLogs, log)
		}
	}

	return filteredLogs, nil
}

// GetAuditLogsByAction filters audit logs by action
func (s *AuditServiceRedis) GetAuditLogsByAction(action string) ([]models.AuditLog, error) {
	allLogs, err := s.GetAllAuditLogs()
	if err != nil {
		return nil, err
	}

	filteredLogs := make([]models.AuditLog, 0)
	for _, log := range allLogs {
		if log.Action == action {
			filteredLogs = append(filteredLogs, log)
		}
	}

	return filteredLogs, nil
}
