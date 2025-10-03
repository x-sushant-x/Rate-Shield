package service

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/config"
	"github.com/x-sushant-x/RateShield/models"
)

// ConfigRulesService implements RulesService interface for file-based configuration
type ConfigRulesService struct {
	configLoader *config.ConfigLoader
	rules        map[string]*models.Rule
}

// NewConfigRulesService creates a new config-based rules service
func NewConfigRulesService(configPath string) (*ConfigRulesService, error) {
	loader := config.NewConfigLoader(configPath)
	
	if err := loader.LoadRules(); err != nil {
		return nil, err
	}

	return &ConfigRulesService{
		configLoader: loader,
		rules:        loader.GetRules(),
	}, nil
}

// GetAllRules returns all rules from the configuration
func (s *ConfigRulesService) GetAllRules() ([]models.Rule, error) {
	rules := make([]models.Rule, 0, len(s.rules))
	
	for _, rule := range s.rules {
		rules = append(rules, *rule)
	}
	
	return rules, nil
}

// GetPaginatedRules returns paginated rules
func (s *ConfigRulesService) GetPaginatedRules(page, items int) (models.PaginatedRules, error) {
	allRules, err := s.GetAllRules()
	if err != nil {
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

// GetRule returns a specific rule by key
func (s *ConfigRulesService) GetRule(key string) (*models.Rule, bool, error) {
	rule, exists := s.rules[key]
	if !exists {
		return nil, false, nil
	}
	return rule, true, nil
}

// SearchRule searches for rules containing the search text in their endpoint
func (s *ConfigRulesService) SearchRule(searchText string) ([]models.Rule, error) {
	var searchedRules []models.Rule
	
	for _, rule := range s.rules {
		if strings.Contains(rule.APIEndpoint, searchText) {
			searchedRules = append(searchedRules, *rule)
		}
	}
	
	return searchedRules, nil
}

// CreateOrUpdateRule is not supported for config-based rules (read-only)
func (s *ConfigRulesService) CreateOrUpdateRule(rule models.Rule) error {
	return errors.New("creating/updating rules is not supported for config-based rules - please modify the configuration file")
}

// DeleteRule is not supported for config-based rules (read-only)
func (s *ConfigRulesService) DeleteRule(endpoint string) error {
	return errors.New("deleting rules is not supported for config-based rules - please modify the configuration file")
}

// CacheRulesLocally returns a copy of the rules map for local caching
func (s *ConfigRulesService) CacheRulesLocally() *map[string]*models.Rule {
	// Create a copy of the rules map
	cachedRules := make(map[string]*models.Rule)
	for key, rule := range s.rules {
		cachedRules[key] = rule
	}
	
	log.Info().Msg("Rules locally cached from config file ✅")
	return &cachedRules
}

// ListenToRulesUpdate is not applicable for config-based rules
// In a real implementation, you might want to watch the file for changes
func (s *ConfigRulesService) ListenToRulesUpdate(updatesChannel chan string) {
	// For config-based rules, we don't have real-time updates
	// This could be enhanced to watch file changes using fsnotify
	log.Info().Msg("File-based rules don't support real-time updates. Restart the application to reload configuration.")
}

// ReloadRules reloads rules from the configuration file
func (s *ConfigRulesService) ReloadRules() error {
	if err := s.configLoader.LoadRules(); err != nil {
		return err
	}
	
	s.rules = s.configLoader.GetRules()
	log.Info().Msg("Rules reloaded from configuration file")
	return nil
}
