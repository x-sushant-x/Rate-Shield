package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/x-sushant-x/RateShield/config"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/service"
	"gopkg.in/yaml.v3"
)

// Migration script to export Redis rules to configuration files
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run migrate_redis_to_config.go <output_format> <output_file>")
		fmt.Println("Example: go run migrate_redis_to_config.go yaml rules_config.yaml")
		fmt.Println("Example: go run migrate_redis_to_config.go json rules_config.json")
		os.Exit(1)
	}

	format := os.Args[1]
	outputFile := os.Args[2]

	if format != "yaml" && format != "json" {
		fmt.Println("❌ Invalid format. Supported formats: yaml, json")
		os.Exit(1)
	}

	fmt.Println("🔄 Connecting to Redis...")
	
	// Connect to Redis
	redisRulesClient, err := redisClient.NewRulesClient()
	if err != nil {
		fmt.Printf("❌ Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}

	redisRulesSvc := service.NewRedisRulesService(redisRulesClient)

	fmt.Println("📥 Fetching rules from Redis...")
	
	// Get all rules from Redis
	rules, err := redisRulesSvc.GetAllRules()
	if err != nil {
		fmt.Printf("❌ Failed to fetch rules from Redis: %v\n", err)
		os.Exit(1)
	}

	if len(rules) == 0 {
		fmt.Println("⚠️  No rules found in Redis")
		os.Exit(0)
	}

	fmt.Printf("📊 Found %d rules in Redis\n", len(rules))

	// Create config structure
	rulesConfig := config.RulesConfig{
		Rules: rules,
	}

	// Marshal to appropriate format
	var data []byte
	switch format {
	case "yaml":
		data, err = yaml.Marshal(rulesConfig)
		if err != nil {
			fmt.Printf("❌ Failed to marshal to YAML: %v\n", err)
			os.Exit(1)
		}
	case "json":
		data, err = json.MarshalIndent(rulesConfig, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to marshal to JSON: %v\n", err)
			os.Exit(1)
		}
	}

	// Write to file
	err = os.WriteFile(outputFile, data, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to write to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully exported %d rules to %s\n", len(rules), outputFile)
	fmt.Println("🔧 You can now use this configuration file by placing it in your application directory")
	fmt.Println("📝 Review the generated file and restart your application to use file-based configuration")

	// Display summary
	fmt.Printf("\n📋 Rules Summary:\n")
	strategyCount := make(map[string]int)
	for _, rule := range rules {
		strategyCount[rule.Strategy]++
		fmt.Printf("   • %s (%s) - %s\n", rule.APIEndpoint, rule.HTTPMethod, rule.Strategy)
	}

	fmt.Printf("\n📊 Strategy Distribution:\n")
	for strategy, count := range strategyCount {
		fmt.Printf("   • %s: %d rules\n", strategy, count)
	}
}
