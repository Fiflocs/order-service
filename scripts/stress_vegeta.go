package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	url := "http://localhost:8080/orders/ORD-2024-001"
	var success, errors int64
	var totalTime int64
	var wg sync.WaitGroup

	requests := 30
	workers := 5

	start := time.Now()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			client := &http.Client{Timeout: 5 * time.Second}

			for j := 0; j < requests/workers; j++ {
				reqStart := time.Now()
				resp, err := client.Get(url)
				latency := time.Since(reqStart).Milliseconds()

				atomic.AddInt64(&totalTime, latency)

				if err == nil && resp.StatusCode == 200 {
					atomic.AddInt64(&success, 1)
					resp.Body.Close()
				} else {
					atomic.AddInt64(&errors, 1)
					fmt.Printf("Error: %v\n", err)
				}

				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start).Seconds()

	fmt.Printf("=== STRESS TEST RESULTS ===\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total requests: %d\n", requests)
	fmt.Printf("Successful: %d\n", success)
	fmt.Printf("Errors: %d\n", errors)
	fmt.Printf("Duration: %.2f seconds\n", duration)
	fmt.Printf("RPS: %.2f\n", float64(success)/duration)
	fmt.Printf("Average latency: %.2f ms\n", float64(totalTime)/float64(success))
}
