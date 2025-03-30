package tests

import (
	"context"
	"errors"
	"ethereum-validator-api/service"
	"ethereum-validator-api/utils"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Helper function to add delay between API calls
func waitForRateLimit() {
	time.Sleep(1100 * time.Millisecond) // Wait slightly more than 1 second to respect rate limit
}

// isRateLimitError checks if the error is related to rate limiting
func isRateLimitError(err error) bool {
	return err != nil && (
		strings.Contains(err.Error(), "429 Too Many Requests") ||
		strings.Contains(err.Error(), "request limit reached") ||
		strings.Contains(err.Error(), "rate limit"))
}

// isAcceptableError checks if an error can be considered acceptable for test cases
// where we don't strictly expect an error but can tolerate certain types of errors
func isAcceptableError(err error) bool {
	return errors.Is(err, service.ErrSlotNotFound)
}

func TestEthereumService_GetBlockRewardBySlot(t *testing.T) {
	// Set up environment variable for testing
	utils.InitializeENV(".env")
	rpcUrl := os.Getenv("ETH_RPC")
	if rpcUrl == "" {
		t.Skip("ETH_RPC environment variable not set, skipping test")
	}
	
	ethService, err := service.NewEthereumService(rpcUrl)
	if err != nil {
		t.Fatalf("Failed to create EthereumService: %v", err)
	}

	// Helper function to check if error is due to rate limiting
	isRateLimitError := func(err error) bool {
		return err != nil && (strings.Contains(err.Error(), "rate limit") ||
			strings.Contains(err.Error(), "too many requests") ||
			strings.Contains(err.Error(), "429"))
	}

	// Helper function to check if the error is acceptable in some test contexts
	isAcceptableError := func(err error) bool {
		return errors.Is(err, service.ErrSlotNotFound)
	}

	// Calculate a recent slot for more accurate testing
	currentSlot := time.Now().Unix() / 12
	recentSlot := currentSlot - 100 // Around 20 minutes ago

	tests := []struct {
		name          string
		slot          int64
		wantErr       bool
		expectedError error
		acceptableErr bool // Flag to indicate if certain errors like ErrSlotNotFound are acceptable
	}{
		{
			name:          "Recent valid slot",
			slot:          recentSlot,
			wantErr:       false,
			acceptableErr: true, // It's acceptable if a recent slot has no data
		},
		{
			name:          "Future slot",
			slot:          currentSlot + 1000,
			wantErr:       true,
			expectedError: service.ErrFutureSlot,
			acceptableErr: false,
		},
		{
			name:          "Very old slot",
			slot:          18900000, // A known valid Ethereum slot
			wantErr:       false,
			acceptableErr: false,
		},
	}

	ctx := context.Background()

	for i, tt := range tests {
		// Add delay between tests to respect rate limit
		if i > 0 {
			waitForRateLimit()
		}

		t.Run(tt.name, func(t *testing.T) {
			var reward *service.BlockReward
			var err error

			// Retry logic for rate limiting
			maxRetries := 3
			for retry := 0; retry < maxRetries; retry++ {
				if retry > 0 {
					waitForRateLimit()
					t.Logf("Retrying %s (attempt %d/%d)", tt.name, retry+1, maxRetries)
				}

				reward, err = ethService.GetBlockRewardBySlot(ctx, tt.slot)
				if err == nil || (tt.wantErr && !isRateLimitError(err)) {
					break
				}

				// If it's not a rate limit error, don't retry
				if err != nil && !isRateLimitError(err) {
					break
				}
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetBlockRewardBySlot() expected error for %s but got nil", tt.name)
					return
				}
				
				// Check for specific error type if expected
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("GetBlockRewardBySlot() expected error %v but got %v", tt.expectedError, err)
				}
				return
			}

			// For non-error cases
			if err != nil {
				// Check if the error is an acceptable one, like ErrSlotNotFound
				if tt.acceptableErr && isAcceptableError(err) {
					t.Logf("Slot %d not found, but this is acceptable for slots that might not exist", tt.slot)
					return
				}
				t.Errorf("GetBlockRewardBySlot() unexpected error: %v", err)
				return
			}

			// Verify the reward structure
			if reward == nil {
				t.Error("Expected non-nil reward")
				return
			}

			// Verify status is either "mev" or "vanilla"
			if reward.Status != "mev" && reward.Status != "vanilla" {
				t.Errorf("Invalid status: %s, expected 'mev' or 'vanilla'", reward.Status)
			}

			// Verify reward is non-negative
			if reward.Reward.Cmp(big.NewInt(0)) == -1 {
				t.Errorf("Expected non-negative reward, got %s", reward.Reward.String())
			}

			t.Logf("Slot %d: Status=%s, Reward=%s GWEI", tt.slot, reward.Status, reward.Reward.String())
		})
	}

	t.Run("Recent valid slot", func(t *testing.T) {
		// Try up to 5 times to get rewards for a recent valid slot
		var reward *service.BlockReward
		var err error
		var slot int64 = time.Now().Unix()/12 - 32 // ~6 minutes ago

		for i := 0; i < 5; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			reward, err = ethService.GetBlockRewardBySlot(ctx, slot)
			if isRateLimitError(err) {
				t.Logf("Rate limited on attempt %d, waiting before retry...", i+1)
				waitForRateLimit()
				continue
			}
			
			if isAcceptableError(err) {
				t.Logf("Slot %d not found, but this is acceptable for recent slots that might not exist", slot)
				return // This is a pass, not a failure
			}
			
			break
		}

		if err != nil {
			t.Fatalf("Failed to get reward for recent slot %d: %v", slot, err)
		}

		// If we got a reward, verify the structure
		if reward != nil {
			t.Logf("Slot %d: Status=%s, Reward=%s GWEI", slot, reward.Status, reward.Reward.String())
			assert.NotNil(t, reward.Status)
			assert.True(t, reward.Reward.Cmp(big.NewInt(0)) >= 0, "Reward should be non-negative")
			assert.True(t, reward.Status == "mev" || reward.Status == "vanilla", "Status should be mev or vanilla")
		}
	})
}

func TestEthereumService_GetSyncDutiesBySlot(t *testing.T) {
	// Set up environment variable for testing
	utils.InitializeENV(".env")
	rpcUrl := os.Getenv("ETH_RPC")
	if rpcUrl == "" {
		t.Skip("ETH_RPC environment variable not set, skipping test")
	}
	
	ethService, err := service.NewEthereumService(rpcUrl)
	if err != nil {
		t.Fatalf("Failed to create EthereumService: %v", err)
	}

	// Calculate a recent slot
	currentSlot := time.Now().Unix() / 12
	recentSlot := currentSlot - (60 * 5) // 5 minutes ago

	tests := []struct {
		name          string
		slot          int64
		wantErr       bool
		expectedError error
	}{
		{
			name:    "Recent valid slot",
			slot:    recentSlot,
			wantErr: false,
		},
		{
			name:          "Future slot",
			slot:          currentSlot + 1000,
			wantErr:       true,
			expectedError: service.ErrFutureSlot,
		},
		{
			name:    "Known historical slot",
			slot:    18900000, // A known Ethereum slot
			wantErr: false,
		},
	}

	ctx := context.Background()

	for i, tt := range tests {
		// Add delay between tests to respect rate limit
		if i > 0 {
			waitForRateLimit()
		}

		t.Run(tt.name, func(t *testing.T) {
			var duties []string
			var err error

			// Retry logic for rate limiting
			maxRetries := 3
			for retry := 0; retry < maxRetries; retry++ {
				if retry > 0 {
					waitForRateLimit()
					t.Logf("Retrying %s (attempt %d/%d)", tt.name, retry+1, maxRetries)
				}

				duties, err = ethService.GetSyncDutiesBySlot(ctx, tt.slot)
				if err == nil || (tt.wantErr && !isRateLimitError(err)) {
					break
				}

				// If it's not a rate limit error, don't retry
				if err != nil && !isRateLimitError(err) {
					break
				}
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetSyncDutiesBySlot() expected error for %s but got nil", tt.name)
					return
				}
				
				// Check for specific error type if expected
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("GetSyncDutiesBySlot() expected error %v but got %v", tt.expectedError, err)
				}
				return
			}

			// For non-error cases
			if err != nil {
				t.Errorf("GetSyncDutiesBySlot() unexpected error: %v", err)
				return
			}

			// Log the number of validators found
			t.Logf("Slot %d: Found %d validator duties", tt.slot, len(duties))

			// Verify that we got some validator public keys
			if len(duties) == 0 {
				t.Logf("Warning: No validator duties found for slot %d", tt.slot)
			} else {
				for i, pubKey := range duties {
					if len(pubKey) == 0 {
						t.Errorf("Empty public key found at index %d", i)
					}
				}
			}
		})
	}
}
