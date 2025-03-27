package tests

import (
	"ethereum-validator-api/utils"
	"os"
	"testing"
)

func TestInitializeENV(t *testing.T) {
	// Create a temporary test .env file
	testEnvContent := []byte("ETH_RPC=https://test-endpoint.example.com")
	err := os.WriteFile("test.env", testEnvContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	// Clean up the test .env file after the test
	defer func() {
		err := os.Remove("test.env")
		if err != nil {
			t.Logf("Failed to remove test .env file: %v", err)
		}
	}()

	tests := []struct {
		name     string
		envFile  string
		wantBool bool
	}{
		{
			name:     "Existing env file",
			envFile:  "test.env",
			wantBool: true,
		},
		{
			name:     "Non-existent env file",
			envFile:  "non_existent.env",
			wantBool: false,
		},
		{
			name:     "Empty filename",
			envFile:  "",
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear any existing environment variables before each test
			err := os.Unsetenv("ETH_RPC")
			if err != nil {
				return
			}

			got := utils.InitializeENV(tt.envFile)
			if got != tt.wantBool {
				t.Errorf("InitializeENV() = %v, want %v", got, tt.wantBool)
			}

			// For successful cases, verify that environment variables were actually loaded
			if tt.wantBool {
				// Check if environment variables are set correctly
				ethRPC := os.Getenv("ETH_RPC")
				if ethRPC != "https://test-endpoint.example.com" {
					t.Errorf("ETH_RPC not set correctly, got: %s, want: %s", ethRPC, "https://test-endpoint.example.com")
				}
			}
		})
	}
}

func TestInitializeENV_InvalidPermissions(t *testing.T) { // Create a test .env file with no read permissions
	testEnvContent := []byte("ETH_RPC=https://test-endpoint.example.com")
	err := os.WriteFile("no_permission.env", testEnvContent, 0000)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	defer func() {
		err := os.Remove("no_permission.env")
		if err != nil {
			t.Logf("Failed to remove test .env file: %v", err)
		}
	}()

	// Test with no-permission file
	got := utils.InitializeENV("no_permission.env")
	if got {
		t.Error("InitializeENV() = true, want false for file with no permissions")
	}
}
