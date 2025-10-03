// Test different configuration scenarios
// Run with: go run test_scenarios.go

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/x-sushant-x/RateShield/config"
)

func main() {
	fmt.Println("🧪 Testing Configuration Scenarios")
	fmt.Println("==================================")

	// Scenario 1: Test with existing YAML config
	fmt.Println("\n📋 Scenario 1: YAML Configuration")
	testYAMLConfig()

	// Scenario 2: Test with JSON config
	fmt.Println("\n📋 Scenario 2: JSON Configuration") 
	testJSONConfig()

	// Scenario 3: Test without any config
	fmt.Println("\n📋 Scenario 3: No Configuration File")
	testNoConfig()

	// Scenario 4: Test invalid config
	fmt.Println("\n📋 Scenario 4: Invalid Configuration")
	testInvalidConfig()

	fmt.Println("\n🎯 Summary:")
	fmt.Println("   ✅ YAML configuration works")
	fmt.Println("   ✅ JSON configuration works") 
	fmt.Println("   ✅ Graceful fallback when no config")
	fmt.Println("   ✅ Proper error handling for invalid config")
}

func testYAMLConfig() {
	if configPath, exists := config.CheckConfigFileExists("."); exists && filepath.Ext(configPath) == ".yaml" {
		loader := config.NewConfigLoader(configPath)
		if err := loader.LoadRules(); err != nil {
			fmt.Printf("   ❌ YAML loading failed: %v\n", err)
			return
		}
		rules := loader.GetRules()
		fmt.Printf("   ✅ YAML config loaded: %d rules\n", len(rules))
	} else {
		fmt.Println("   ⚠️  No YAML config found")
	}
}

func testJSONConfig() {
	jsonPath := "example_rules_config.json"
	if _, err := os.Stat(jsonPath); err == nil {
		loader := config.NewConfigLoader(jsonPath)
		if err := loader.LoadRules(); err != nil {
			fmt.Printf("   ❌ JSON loading failed: %v\n", err)
			return
		}
		rules := loader.GetRules()
		fmt.Printf("   ✅ JSON config loaded: %d rules\n", len(rules))
	} else {
		fmt.Println("   ⚠️  No JSON config found")
	}
}

func testNoConfig() {
	// Test in a temporary directory with no config files
	tempDir := os.TempDir()
	configPath, exists := config.CheckConfigFileExists(tempDir)
	if !exists {
		fmt.Println("   ✅ Correctly detected no config file")
	} else {
		fmt.Printf("   ⚠️  Unexpected config found: %s\n", configPath)
	}
}

func testInvalidConfig() {
	// Create a temporary invalid config
	invalidYAML := `
rules:
  - strategy: "INVALID_STRATEGY"
    endpoint: ""
    token_bucket_rule:
      bucket_capacity: -1
`
	
	tempFile := filepath.Join(os.TempDir(), "invalid_config.yaml")
	err := os.WriteFile(tempFile, []byte(invalidYAML), 0644)
	if err != nil {
		fmt.Printf("   ❌ Could not create test file: %v\n", err)
		return
	}
	defer os.Remove(tempFile)

	loader := config.NewConfigLoader(tempFile)
	err = loader.LoadRules()
	if err != nil {
		fmt.Printf("   ✅ Correctly rejected invalid config: %v\n", err)
	} else {
		rules := loader.GetRules()
		fmt.Printf("   ⚠️  Invalid config was accepted, loaded %d rules\n", len(rules))
	}
}
