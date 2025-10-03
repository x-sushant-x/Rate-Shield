package main

import (
	"fmt"
	"os"

	"github.com/x-sushant-x/RateShield/config"
)

// Simple validation script to test configuration files
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run validate_config.go <config_file_path>")
		fmt.Println("Example: go run validate_config.go rules_config.yaml")
		os.Exit(1)
	}

	configPath := os.Args[1]
	
	fmt.Printf("Validating configuration file: %s\n", configPath)
	
	loader := config.NewConfigLoader(configPath)
	err := loader.LoadRules()
	if err != nil {
		fmt.Printf("❌ Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	rules := loader.GetRules()
	fmt.Printf("✅ Configuration is valid!\n")
	fmt.Printf("📊 Loaded %d rules:\n\n", len(rules))

	for endpoint, rule := range rules {
		fmt.Printf("🔧 Endpoint: %s\n", endpoint)
		fmt.Printf("   Strategy: %s\n", rule.Strategy)
		fmt.Printf("   HTTP Method: %s\n", rule.HTTPMethod)
		fmt.Printf("   Allow on Error: %t\n", rule.AllowOnError)
		
		switch rule.Strategy {
		case "TOKEN BUCKET":
			if rule.TokenBucketRule != nil {
				fmt.Printf("   Bucket Capacity: %d\n", rule.TokenBucketRule.BucketCapacity)
				fmt.Printf("   Token Add Rate: %d\n", rule.TokenBucketRule.TokenAddRate)
				fmt.Printf("   Retention Time: %d seconds\n", rule.TokenBucketRule.RetentionTime)
			}
		case "FIXED WINDOW COUNTER":
			if rule.FixedWindowCounterRule != nil {
				fmt.Printf("   Max Requests: %d\n", rule.FixedWindowCounterRule.MaxRequests)
				fmt.Printf("   Window: %d seconds\n", rule.FixedWindowCounterRule.Window)
			}
		case "SLIDING WINDOW COUNTER":
			if rule.SlidingWindowCounterRule != nil {
				fmt.Printf("   Max Requests: %d\n", rule.SlidingWindowCounterRule.MaxRequests)
				fmt.Printf("   Window Size: %d seconds\n", rule.SlidingWindowCounterRule.WindowSize)
			}
		}
		fmt.Println()
	}
}
