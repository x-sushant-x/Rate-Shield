package service

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
)

type RulesService interface {
	GetAllRules() ([]models.Rule, error)
	GetRule(key string) (*models.Rule, bool, error)
	SearchRule(searchText string) ([]models.Rule, error)
	CreateOrUpdateRule(models.Rule) error
	DeleteRule(endpoint string) error
}

type RulesServiceRedis struct {
	redisClient redisClient.RedisRuleClient
}

func NewRedisRulesService(client redisClient.RedisRuleClient) RulesServiceRedis {
	return RulesServiceRedis{
		redisClient: client,
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

func (s RulesServiceRedis) CreateOrUpdateRule(rule models.Rule) error {
	err := s.redisClient.SetRule(rule.APIEndpoint, rule)
	if err != nil {
		log.Err(err).Msg("unable to create or update rule")
		return err
	}
	return nil
}

func (s RulesServiceRedis) DeleteRule(endpoint string) error {
	err := s.redisClient.DeleteRule(endpoint)
	if err != nil {
		log.Err(err).Msg("unable to create or update rule")
		return err

	}
	return err
}
