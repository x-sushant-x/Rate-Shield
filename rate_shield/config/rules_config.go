package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	"gopkg.in/yaml.v3"
)

// RulesConfig represents the structure of the configuration file
type RulesConfig struct {
	Rules []models.Rule `yaml:"rules" json:"rules"`
}

// ConfigLoader handles loading rules from YAML/JSON configuration files
type ConfigLoader struct {
	configPath string
	rules      map[string]*models.Rule
}

// NewConfigLoader creates a new config loader instance
func NewConfigLoader(configPath string) *ConfigLoader {
	return &ConfigLoader{
		configPath: configPath,
		rules:      make(map[string]*models.Rule),
	}
}

// LoadRules loads rules from the configuration file
func (c *ConfigLoader) LoadRules() error {
	// Check if file exists
	if _, err := os.Stat(c.configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", c.configPath)
	}

	// Read file content
	data, err := os.ReadFile(c.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse based on file extension
	var rulesConfig RulesConfig
	ext := strings.ToLower(filepath.Ext(c.configPath))

	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &rulesConfig)
	case ".json":
		err = json.Unmarshal(data, &rulesConfig)
	default:
		return fmt.Errorf("unsupported file format: %s. Supported formats: .yaml, .yml, .json", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and store rules
	c.rules = make(map[string]*models.Rule)
	for i, rule := range rulesConfig.Rules {
		if err := c.validateRule(&rule); err != nil {
			log.Warn().Msgf("Skipping invalid rule at index %d: %v", i, err)
			continue
		}
		c.rules[rule.APIEndpoint] = &rule
	}

	log.Info().Msgf("Loaded %d rules from config file: %s", len(c.rules), c.configPath)
	return nil
}

// validateRule validates a rule configuration
func (c *ConfigLoader) validateRule(rule *models.Rule) error {
	if rule.APIEndpoint == "" {
		return fmt.Errorf("API endpoint is required")
	}

	if rule.Strategy == "" {
		return fmt.Errorf("strategy is required")
	}

	// Validate strategy-specific rules
	switch rule.Strategy {
	case "TOKEN BUCKET":
		if rule.TokenBucketRule == nil {
			return fmt.Errorf("token bucket rule is required for TOKEN BUCKET strategy")
		}
		if rule.TokenBucketRule.BucketCapacity <= 0 {
			return fmt.Errorf("bucket capacity must be positive")
		}
		if rule.TokenBucketRule.TokenAddRate <= 0 {
			return fmt.Errorf("token add rate must be positive")
		}
	case "FIXED WINDOW COUNTER":
		if rule.FixedWindowCounterRule == nil {
			return fmt.Errorf("fixed window counter rule is required for FIXED WINDOW COUNTER strategy")
		}
		if rule.FixedWindowCounterRule.MaxRequests <= 0 {
			return fmt.Errorf("max requests must be positive")
		}
		if rule.FixedWindowCounterRule.Window <= 0 {
			return fmt.Errorf("window must be positive")
		}
	case "SLIDING WINDOW COUNTER":
		if rule.SlidingWindowCounterRule == nil {
			return fmt.Errorf("sliding window counter rule is required for SLIDING WINDOW COUNTER strategy")
		}
		if rule.SlidingWindowCounterRule.MaxRequests <= 0 {
			return fmt.Errorf("max requests must be positive")
		}
		if rule.SlidingWindowCounterRule.WindowSize <= 0 {
			return fmt.Errorf("window size must be positive")
		}
	default:
		return fmt.Errorf("unsupported strategy: %s", rule.Strategy)
	}

	return nil
}

// GetRules returns all loaded rules
func (c *ConfigLoader) GetRules() map[string]*models.Rule {
	return c.rules
}

// GetRule returns a specific rule by endpoint
func (c *ConfigLoader) GetRule(endpoint string) (*models.Rule, bool) {
	rule, exists := c.rules[endpoint]
	return rule, exists
}

// CheckConfigFileExists checks if any of the supported config files exist
func CheckConfigFileExists(basePath string) (string, bool) {
	configFiles := []string{
		"rules_config.yaml",
		"rules_config.yml", 
		"rules_config.json",
	}

	for _, filename := range configFiles {
		fullPath := filepath.Join(basePath, filename)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, true
		}
	}

	return "", false
}
