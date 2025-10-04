package service

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
)

const (
	redisChannel = "rules-update"
)

type RulesService interface {
	GetAllRules() ([]models.Rule, error)
	GetPaginatedRules(page, items int) (models.PaginatedRules, error)
	GetRule(key string) (*models.Rule, bool, error)
	SearchRule(searchText string) ([]models.Rule, error)
	CreateOrUpdateRule(rule models.Rule, actor, ipAddress, userAgent string) error
	DeleteRule(endpoint, actor, ipAddress, userAgent string) error
	CacheRulesLocally() *map[string]*models.Rule
	ListenToRulesUpdate(updatesChannel chan string)
}

type RulesServiceRedis struct {
	redisClient redisClient.RedisRuleClient
	auditSvc    AuditService
}

func NewRedisRulesService(client redisClient.RedisRuleClient, auditSvc AuditService) RulesServiceRedis {
	return RulesServiceRedis{
		redisClient: client,
		auditSvc:    auditSvc,
	}
}

func (s RulesServiceRedis) GetRule(key string) (*models.Rule, bool, error) {
	return s.redisClient.GetRule(key)
}

func (s RulesServiceRedis) GetAllRules() ([]models.Rule, error) {
	keys, _, err := s.redisClient.GetAllRuleKeys()
	if err != nil {
		log.Err(err).Msg("unable to get all rule keys from redis")
	}

	rules := []models.Rule{}

	for _, key := range keys {
		rule, found, err := s.redisClient.GetRule(key)

		if !found {
			log.Error().Msgf("rule with key: %s not found", key)
			continue
		}

		if err != nil {
			log.Err(err).Msg("unable to get rule from redis")
			continue
		}

		rules = append(rules, *rule)
	}

	return rules, nil
}

func (s RulesServiceRedis) SearchRule(searchText string) ([]models.Rule, error) {
	rules, err := s.GetAllRules()
	if err != nil {
		return nil, err
	}
	searchedRules := []models.Rule{}

	for _, rule := range rules {
		if strings.Contains(rule.APIEndpoint, searchText) {
			searchedRules = append(searchedRules, rule)
		}
	}

	return searchedRules, nil
}

func (s RulesServiceRedis) CreateOrUpdateRule(rule models.Rule, actor, ipAddress, userAgent string) error {
	// Check if rule already exists to determine action (CREATE vs UPDATE)
	existingRule, found, err := s.redisClient.GetRule(rule.APIEndpoint)

	var action string
	var oldRule *models.Rule

	if found && err == nil {
		// Rule exists - this is an UPDATE
		action = models.AuditActionUpdate
		oldRule = existingRule
	} else {
		// Rule doesn't exist - this is a CREATE
		action = models.AuditActionCreate
		oldRule = nil
	}

	// Save the rule to Redis
	err = s.redisClient.SetRule(rule.APIEndpoint, rule)
	if err != nil {
		log.Err(err).Msg("unable to create or update rule")
		return err
	}

	// Log audit event
	if s.auditSvc != nil {
		auditErr := s.auditSvc.LogRuleChange(actor, action, rule.APIEndpoint, oldRule, &rule, ipAddress, userAgent)
		if auditErr != nil {
			log.Warn().Err(auditErr).Msg("failed to log audit event for rule change")
			// Don't fail the operation if audit logging fails
		}
	}

	return s.redisClient.PublishMessage(redisChannel, "rule-updated")
}

func (s RulesServiceRedis) DeleteRule(endpoint, actor, ipAddress, userAgent string) error {
	// Get the existing rule before deleting for audit log
	existingRule, found, err := s.redisClient.GetRule(endpoint)
	if !found || err != nil {
		log.Warn().Str("endpoint", endpoint).Msg("rule not found for deletion")
		// Still attempt to delete in case of inconsistency
	}

	// Delete the rule from Redis
	err = s.redisClient.DeleteRule(endpoint)
	if err != nil {
		log.Err(err).Msg("unable to delete rule")
		return err
	}

	// Log audit event
	if s.auditSvc != nil && existingRule != nil {
		auditErr := s.auditSvc.LogRuleChange(actor, models.AuditActionDelete, endpoint, existingRule, nil, ipAddress, userAgent)
		if auditErr != nil {
			log.Warn().Err(auditErr).Msg("failed to log audit event for rule deletion")
			// Don't fail the operation if audit logging fails
		}
	}

	return s.redisClient.PublishMessage(redisChannel, "rule-updated")
}

func (s RulesServiceRedis) CacheRulesLocally() *map[string]*models.Rule {
	rules, err := s.GetAllRules()
	if err != nil {
		log.Err(err).Msg("Unable to cache all rules locally")
	}

	cachedRules := make(map[string]*models.Rule)

	for _, rule := range rules {
		cachedRules[rule.APIEndpoint] = &rule
	}

	log.Info().Msg("Rules locally cached ✅")
	return &cachedRules
}

func (s RulesServiceRedis) ListenToRulesUpdate(updatesChannel chan string) {
	s.redisClient.ListenToRulesUpdate(updatesChannel)
}

func (s RulesServiceRedis) GetPaginatedRules(page, items int) (models.PaginatedRules, error) {
	allRules, err := s.GetAllRules()
	if err != nil {
		log.Err(err).Msgf("unable to get rules from redis")
		return models.PaginatedRules{}, err
	}

	if len(allRules) == 0 {
		return models.PaginatedRules{
			PageNumber:  1,
			TotalItems:  0,
			HasNextPage: false,
			Rules:       make([]models.Rule, 0),
		}, nil
	}

	start := (page - 1) * items
	stop := start + items

	if start >= len(allRules) {
		return models.PaginatedRules{}, errors.New("invalid page number")
	}

	hasNextPage := stop < len(allRules)

	if stop >= len(allRules) {
		stop = len(allRules)
	}

	paginatedSlice := allRules[start:stop]

	rules := models.PaginatedRules{
		PageNumber:  page,
		TotalItems:  stop - start,
		HasNextPage: hasNextPage,
		Rules:       paginatedSlice,
	}

	return rules, nil
}
