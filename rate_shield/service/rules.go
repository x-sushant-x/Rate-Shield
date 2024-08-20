package service

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
)

type RulesService interface {
	GetAllRules() ([]models.Rule, error)
	UpdateRule(endpoint string) error
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

		var r models.Rule
		err = json.Unmarshal(rule, &r)
		if err != nil {
			log.Err(err).Msgf("unable to marshal rule with key: %s", key)
			continue
		}

		rules = append(rules, r)
	}

	return rules, nil
}

func (s RulesServiceRedis) UpdateRule(endpoint string) error { return nil }
func (s RulesServiceRedis) DeleteRule(endpoint string) error { return nil }
