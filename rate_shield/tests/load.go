package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	totalRequests  = 500000
	maxConcurrency = 100

	successResponse         = 0
	tooManyRequestsResponse = 0
)

func main() {
	var IPs []net.IP

	for i := 0; i < 1000; i++ {
		IPs = append(IPs, generateRandomIP())
	}

	startTime := time.Now()

	var wg sync.WaitGroup

	semaphore := make(chan struct{}, maxConcurrency)

	var initalMem runtime.MemStats
	runtime.ReadMemStats(&initalMem)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ip := pickRandomIP(IPs)

			req, err := http.NewRequest("GET", "http://127.0.0.1:8080/check-limit", nil)
			if err != nil {
				return
			}

			req.Header.Add("ip", ip.String())
			req.Header.Add("endpoint", "/api/v1/resource")

			res, _ := http.DefaultClient.Do(req)
			if res.StatusCode == 200 {
				successResponse++
			} else if res.StatusCode == 429 {
				tooManyRequestsResponse++
			}
		}()
	}

	wg.Wait()

	var finalMem runtime.MemStats
	runtime.ReadMemStats(&finalMem)

	memUsed := finalMem.Alloc - initalMem.Alloc

	fmt.Printf("Memory Used: %dMB \n", memUsed/(1024*1024))

	duration := time.Since(startTime)
	fmt.Println("Time Taken: ", duration)

	requestsPerSecond := float64(totalRequests) / duration.Seconds()
	fmt.Printf("Requests per second: %.2f\n", requestsPerSecond)

	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Total Success Response: %d\n", successResponse)
	fmt.Printf("Total Too Many Requests Response: %d\n", tooManyRequestsResponse)
}

func generateRandomIP() net.IP {
	ip := net.IPv4(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)))
	return ip
}

func pickRandomIP(IPs []net.IP) net.IP {
	idx := rand.Intn(len(IPs))
	return IPs[idx]
}
