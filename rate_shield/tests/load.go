package main

import (
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	rate := vegeta.Rate{
		Freq: 1000,
		Per:  time.Second,
	}

	duration := time.Second * 60
	attacker := vegeta.NewAttacker()

	ipGenerator := func(i int) string {
		return fmt.Sprintf("192.168.1.%d", i%100+1) // Generates IPs like 192.168.1.1 to 192.168.1.10000
	}

	var metrics vegeta.Metrics
	for i := 0; i < 100; i++ {
		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: "GET",
			URL:    "http://127.0.0.1/check-limit",
			Header: map[string][]string{
				"ip":       {ipGenerator(i)},
				"endpoint": {"/api/v1/get-data"},
			},
		})

		for res := range attacker.Attack(targeter, rate, duration, "Rate Limiter Test") {
			metrics.Add(res)
		}
		metrics.Close()
	}

	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Rate: %f\n", metrics.Rate)
	fmt.Printf("Duration: %s\n", metrics.Duration)
	fmt.Printf("Success: %f\n", metrics.Success)
	fmt.Printf("Latencies: %s\n", metrics.Latencies)

}
