package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x-sushant-x/RateShield/models"
)

func TestConfigLoader_LoadRules_YAML(t *testing.T) {
	// Create a temporary YAML config file
	yamlContent := `
rules:
  - strategy: "TOKEN BUCKET"
    endpoint: "/api/test"
    http_method: "GET"
    allow_on_error: true
    token_bucket_rule:
      bucket_capacity: 100
      token_add_rate: 10
      retention_time: 3600
  - strategy: "FIXED WINDOW COUNTER"
    endpoint: "/api/login"
    http_method: "POST"
    allow_on_error: false
    fixed_window_counter_rule:
      max_requests: 5
      window: 60
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	// Test loading the config
	loader := NewConfigLoader(configPath)
	err = loader.LoadRules()
	assert.NoError(t, err)

	rules := loader.GetRules()
	assert.Len(t, rules, 2)

	// Test token bucket rule
	tokenBucketRule, exists := rules["/api/test"]
	assert.True(t, exists)
	assert.Equal(t, "TOKEN BUCKET", tokenBucketRule.Strategy)
	assert.Equal(t, "/api/test", tokenBucketRule.APIEndpoint)
	assert.True(t, tokenBucketRule.AllowOnError)
	assert.NotNil(t, tokenBucketRule.TokenBucketRule)
	assert.Equal(t, int64(100), tokenBucketRule.TokenBucketRule.BucketCapacity)

	// Test fixed window rule
	fixedWindowRule, exists := rules["/api/login"]
	assert.True(t, exists)
	assert.Equal(t, "FIXED WINDOW COUNTER", fixedWindowRule.Strategy)
	assert.Equal(t, "/api/login", fixedWindowRule.APIEndpoint)
	assert.False(t, fixedWindowRule.AllowOnError)
	assert.NotNil(t, fixedWindowRule.FixedWindowCounterRule)
	assert.Equal(t, int64(5), fixedWindowRule.FixedWindowCounterRule.MaxRequests)
}

func TestConfigLoader_LoadRules_JSON(t *testing.T) {
	// Create a temporary JSON config file
	jsonContent := `{
  "rules": [
    {
      "strategy": "SLIDING WINDOW COUNTER",
      "endpoint": "/api/upload",
      "http_method": "POST",
      "allow_on_error": true,
      "sliding_window_counter_rule": {
        "max_requests": 20,
        "window": 300
      }
    }
  ]
}`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.json")
	
	err := os.WriteFile(configPath, []byte(jsonContent), 0644)
	assert.NoError(t, err)

	// Test loading the config
	loader := NewConfigLoader(configPath)
	err = loader.LoadRules()
	assert.NoError(t, err)

	rules := loader.GetRules()
	assert.Len(t, rules, 1)

	// Test sliding window rule
	slidingWindowRule, exists := rules["/api/upload"]
	assert.True(t, exists)
	assert.Equal(t, "SLIDING WINDOW COUNTER", slidingWindowRule.Strategy)
	assert.Equal(t, "/api/upload", slidingWindowRule.APIEndpoint)
	assert.True(t, slidingWindowRule.AllowOnError)
	assert.NotNil(t, slidingWindowRule.SlidingWindowCounterRule)
	assert.Equal(t, int64(20), slidingWindowRule.SlidingWindowCounterRule.MaxRequests)
}

func TestConfigLoader_InvalidConfig(t *testing.T) {
	// Test invalid YAML
	invalidYAML := `
rules:
  - strategy: "TOKEN BUCKET"
    endpoint: ""  # Invalid: empty endpoint
    token_bucket_rule:
      bucket_capacity: -1  # Invalid: negative capacity
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid_config.yaml")
	
	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	assert.NoError(t, err)

	loader := NewConfigLoader(configPath)
	err = loader.LoadRules()
	assert.NoError(t, err) // Should not error, but should skip invalid rules

	rules := loader.GetRules()
	assert.Len(t, rules, 0) // No valid rules should be loaded
}

func TestCheckConfigFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Test when no config file exists
	configPath, exists := CheckConfigFileExists(tmpDir)
	assert.False(t, exists)
	assert.Empty(t, configPath)

	// Create a YAML config file
	yamlPath := filepath.Join(tmpDir, "rules_config.yaml")
	err := os.WriteFile(yamlPath, []byte("rules: []"), 0644)
	assert.NoError(t, err)

	// Test when YAML config exists
	configPath, exists = CheckConfigFileExists(tmpDir)
	assert.True(t, exists)
	assert.Equal(t, yamlPath, configPath)

	// Create a JSON config file (should still prefer YAML)
	jsonPath := filepath.Join(tmpDir, "rules_config.json")
	err = os.WriteFile(jsonPath, []byte(`{"rules": []}`), 0644)
	assert.NoError(t, err)

	configPath, exists = CheckConfigFileExists(tmpDir)
	assert.True(t, exists)
	assert.Equal(t, yamlPath, configPath) // Should still prefer YAML
}

func TestValidateRule(t *testing.T) {
	loader := NewConfigLoader("")

	// Test valid token bucket rule
	validRule := &models.Rule{
		Strategy:    "TOKEN BUCKET",
		APIEndpoint: "/api/test",
		TokenBucketRule: &models.TokenBucketRule{
			BucketCapacity: 100,
			TokenAddRate:   10,
			RetentionTime:  3600,
		},
	}
	err := loader.validateRule(validRule)
	assert.NoError(t, err)

	// Test invalid rule - missing endpoint
	invalidRule := &models.Rule{
		Strategy: "TOKEN BUCKET",
		TokenBucketRule: &models.TokenBucketRule{
			BucketCapacity: 100,
			TokenAddRate:   10,
			RetentionTime:  3600,
		},
	}
	err = loader.validateRule(invalidRule)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API endpoint is required")

	// Test invalid rule - missing strategy
	invalidRule2 := &models.Rule{
		APIEndpoint: "/api/test",
		TokenBucketRule: &models.TokenBucketRule{
			BucketCapacity: 100,
			TokenAddRate:   10,
			RetentionTime:  3600,
		},
	}
	err = loader.validateRule(invalidRule2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "strategy is required")
}
