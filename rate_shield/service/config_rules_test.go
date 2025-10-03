package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigRulesService_Integration(t *testing.T) {
	// Create a temporary config file
	yamlContent := `
rules:
  - strategy: "TOKEN BUCKET"
    endpoint: "/api/users"
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
  - strategy: "SLIDING WINDOW COUNTER"
    endpoint: "/api/upload"
    http_method: "POST"
    allow_on_error: true
    sliding_window_counter_rule:
      max_requests: 20
      window: 300
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	// Create service
	service, err := NewConfigRulesService(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test GetAllRules
	allRules, err := service.GetAllRules()
	assert.NoError(t, err)
	assert.Len(t, allRules, 3)

	// Test GetRule
	rule, found, err := service.GetRule("/api/users")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.NotNil(t, rule)
	assert.Equal(t, "TOKEN BUCKET", rule.Strategy)
	assert.Equal(t, "/api/users", rule.APIEndpoint)

	// Test GetRule - not found
	rule, found, err = service.GetRule("/api/nonexistent")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, rule)

	// Test SearchRule
	searchResults, err := service.SearchRule("api")
	assert.NoError(t, err)
	assert.Len(t, searchResults, 3) // All rules contain "api"

	searchResults, err = service.SearchRule("login")
	assert.NoError(t, err)
	assert.Len(t, searchResults, 1)
	assert.Equal(t, "/api/login", searchResults[0].APIEndpoint)

	// Test GetPaginatedRules
	paginatedRules, err := service.GetPaginatedRules(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, paginatedRules.PageNumber)
	assert.Equal(t, 2, paginatedRules.TotalItems)
	assert.True(t, paginatedRules.HasNextPage)
	assert.Len(t, paginatedRules.Rules, 2)

	// Test CacheRulesLocally
	cachedRules := service.CacheRulesLocally()
	assert.NotNil(t, cachedRules)
	assert.Len(t, *cachedRules, 3)

	// Test read-only operations (should return errors)
	err = service.CreateOrUpdateRule(allRules[0])
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported for config-based rules")

	err = service.DeleteRule("/api/users")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported for config-based rules")
}

func TestConfigRulesService_InvalidConfig(t *testing.T) {
	// Test with non-existent file
	service, err := NewConfigRulesService("/nonexistent/path/config.yaml")
	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "config file not found")
}

func TestConfigRulesService_EmptyRules(t *testing.T) {
	// Create empty config file
	yamlContent := `rules: []`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty_config.yaml")
	
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	// Create service
	service, err := NewConfigRulesService(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test with empty rules
	allRules, err := service.GetAllRules()
	assert.NoError(t, err)
	assert.Len(t, allRules, 0)

	// Test pagination with empty rules
	paginatedRules, err := service.GetPaginatedRules(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, paginatedRules.PageNumber)
	assert.Equal(t, 0, paginatedRules.TotalItems)
	assert.False(t, paginatedRules.HasNextPage)
	assert.Len(t, paginatedRules.Rules, 0)
}
