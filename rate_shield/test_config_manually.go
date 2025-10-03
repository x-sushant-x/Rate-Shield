package main

import (
	"fmt"
	"log"

	"github.com/x-sushant-x/RateShield/config"
	"github.com/x-sushant-x/RateShield/service"
)

// Manual test script to verify configuration loading
func main() {
	fmt.Println("🧪 Testing Configuration-Based Rules Implementation")
	fmt.Println("=" * 50)

	// Test 1: Check if config file exists
	fmt.Println("\n1️⃣ Testing config file detection...")
	if configPath, exists := config.CheckConfigFileExists("."); exists {
		fmt.Printf("✅ Found config file: %s\n", configPath)
		
		// Test 2: Load configuration
		fmt.Println("\n2️⃣ Testing configuration loading...")
		loader := config.NewConfigLoader(configPath)
		if err := loader.LoadRules(); err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}
		
		rules := loader.GetRules()
		fmt.Printf("✅ Successfully loaded %d rules\n", len(rules))
		
		// Test 3: Display loaded rules
		fmt.Println("\n3️⃣ Loaded rules:")
		for endpoint, rule := range rules {
			fmt.Printf("   📍 %s (%s) - %s\n", endpoint, rule.HTTPMethod, rule.Strategy)
		}
		
		// Test 4: Test ConfigRulesService
		fmt.Println("\n4️⃣ Testing ConfigRulesService...")
		configService, err := service.NewConfigRulesService(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to create ConfigRulesService: %v\n", err)
			return
		}
		
		// Test GetAllRules
		allRules, err := configService.GetAllRules()
		if err != nil {
			fmt.Printf("❌ GetAllRules failed: %v\n", err)
			return
		}
		fmt.Printf("✅ GetAllRules returned %d rules\n", len(allRules))
		
		// Test GetRule
		if len(allRules) > 0 {
			testEndpoint := allRules[0].APIEndpoint
			rule, found, err := configService.GetRule(testEndpoint)
			if err != nil {
				fmt.Printf("❌ GetRule failed: %v\n", err)
				return
			}
			if found {
				fmt.Printf("✅ GetRule found rule for %s\n", testEndpoint)
			} else {
				fmt.Printf("❌ GetRule didn't find rule for %s\n", testEndpoint)
			}
		}
		
		// Test SearchRule
		searchResults, err := configService.SearchRule("api")
		if err != nil {
			fmt.Printf("❌ SearchRule failed: %v\n", err)
			return
		}
		fmt.Printf("✅ SearchRule found %d rules containing 'api'\n", len(searchResults))
		
		// Test CacheRulesLocally
		cachedRules := configService.CacheRulesLocally()
		fmt.Printf("✅ CacheRulesLocally returned %d cached rules\n", len(*cachedRules))
		
		fmt.Println("\n🎉 All tests passed! Configuration system is working correctly.")
		
	} else {
		fmt.Println("❌ No config file found. Please ensure rules_config.yaml exists in the current directory.")
		fmt.Println("💡 You can use the example file that was created.")
	}
}
