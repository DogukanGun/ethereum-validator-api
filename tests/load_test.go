package tests

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestLoadBalancing(t *testing.T) {
	// Configuration
	baseURL := "http://localhost:3004"
	testSlot := "4700000" // Using a known valid slot from the existing tests
	concurrentUsers := 50
	totalRequests := 1000
	duration := "30s"

	tests := []struct {
		name     string
		endpoint string
	}{
		{
			name:     "Block Reward Endpoint",
			endpoint: fmt.Sprintf("%s/blockreward/%s", baseURL, testSlot),
		},
		{
			name:     "Sync Duties Endpoint",
			endpoint: fmt.Sprintf("%s/syncduties/%s", baseURL, testSlot),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build hey command with parameters
			cmd := exec.Command("hey",
				"-n", fmt.Sprintf("%d", totalRequests), // Number of requests
				"-c", fmt.Sprintf("%d", concurrentUsers), // Number of concurrent users
				"-z", duration, // Test duration
				"-q", "5", // Rate limit (requests per second)
				tt.endpoint,
			)

			// Run the command and capture output
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to run hey command: %v\nOutput: %s", err, string(output))
			}

			// Log the results
			t.Logf("Load test results for %s:\n%s", tt.name, string(output))

			// Add delay between tests to allow API to recover
			time.Sleep(5 * time.Second)
		})
	}
}

func TestLoadBalancingWithRateLimit(t *testing.T) {
	// Configuration with rate limiting
	baseURL := "http://localhost:3004"
	testSlot := "4700000"
	concurrentUsers := 20
	totalRequests := 500
	rateLimit := 2 // Requests per second, matching the 2-second delay in the service

	tests := []struct {
		name     string
		endpoint string
	}{
		{
			name:     "Rate Limited Block Reward",
			endpoint: fmt.Sprintf("%s/blockreward/%s", baseURL, testSlot),
		},
		{
			name:     "Rate Limited Sync Duties",
			endpoint: fmt.Sprintf("%s/syncduties/%s", baseURL, testSlot),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build hey command with rate limiting
			cmd := exec.Command("hey",
				"-n", fmt.Sprintf("%d", totalRequests),
				"-c", fmt.Sprintf("%d", concurrentUsers),
				"-q", fmt.Sprintf("%d", rateLimit),
				tt.endpoint,
			)

			// Run the command and capture output
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to run hey command: %v\nOutput: %s", err, string(output))
			}

			// Log the results
			t.Logf("Rate limited load test results for %s:\n%s", tt.name, string(output))

			// Add delay between tests
			time.Sleep(10 * time.Second)
		})
	}
} 