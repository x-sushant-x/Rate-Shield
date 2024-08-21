package service

import (
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
)

type RulesService interface {
	GetAllRules() ([]models.Rule, error)
	CreateOrUpdateRule(models.Rule) error
	DeleteRule(endpoint string) error
}

type RulesServiceRedis struct{}

func (s RulesServiceRedis) GetAllRules() ([]models.Rule, error) {
	keys, _, err := redisClient.GetAllRuleKeys()
	if err != nil {
		log.Err(err).Msg("unable to get all rule keys from redis")
	}

	rules := []models.Rule{}

	for _, key := range keys {
		rule, found, err := redisClient.GetRule(key)

		if !found {
			log.Error().Msgf("rule with key: %s not found", key)
			continue
		}

		if err != nil {
			log.Err(err).Msg("unable to get rule from redis")
			continue
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func (s RulesServiceRedis) CreateOrUpdateRule(rule models.Rule) error {
	err := redisClient.SetRule(rule.APIEndpoint, rule)
	if err != nil {
		log.Err(err).Msg("unable to create or update rule")
		return err
	}
	return nil
}

func (s RulesServiceRedis) DeleteRule(endpoint string) error {
	err := redisClient.DeleteRule(endpoint)
	if err != nil {
		log.Err(err).Msg("unable to create or update rule")
		return err

	}
	return err
}
