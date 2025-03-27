package tests

import (
	"context"
	service2 "ethereum-validator-api/service"
	"ethereum-validator-api/utils"
	"math/big"
	"os"
	_ "os"
	"testing"
	"time"
)

// Helper function to add delay between API calls
func waitForRateLimit() {
	time.Sleep(1100 * time.Millisecond) // Wait slightly more than 1 second to respect rate limit
}

func TestEthereumService_GetBlockRewardBySlot(t *testing.T) {
	// Set up environment variable for testing
	utils.InitializeENV(".env")
	rpcUrl := os.Getenv("ETH_RPC")
	service, err := service2.NewEthereumService(rpcUrl)
	if err != nil {
		t.Fatalf("Failed to create EthereumService: %v", err)
	}

	tests := []struct {
		name    string
		slot    int64
		wantErr bool
	}{
		{
			name:    "Recent valid slot",
			slot:    18900000, // A known Ethereum slot
			wantErr: false,
		},
		{
			name:    "Future slot",
			slot:    999999999,
			wantErr: true,
		},
		{
			name:    "Very old slot",
			slot:    1,
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
			var reward *service2.BlockReward
			var err error

			// Retry logic for rate limiting
			maxRetries := 3
			for retry := 0; retry < maxRetries; retry++ {
				if retry > 0 {
					waitForRateLimit()
					t.Logf("Retrying %s (attempt %d/%d)", tt.name, retry+1, maxRetries)
				}

				reward, err = service.GetBlockRewardBySlot(ctx, tt.slot)
				if err == nil || tt.wantErr {
					break
				}

				// If it's not a rate limit error, don't retry
				if err != nil && err.Error() != "failed to get block: 429 Too Many Requests" {
					break
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockRewardBySlot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
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
			}
		})
	}
}

func TestEthereumService_GetSyncDutiesBySlot(t *testing.T) {
	// Set up environment variable for testing
	utils.InitializeENV(".env")
	rpcUrl := os.Getenv("ETH_RPC")
	service, err := service2.NewEthereumService(rpcUrl)
	if err != nil {
		t.Fatalf("Failed to create EthereumService: %v", err)
	}

	// Calculate a recent slot (about 1 hour ago)
	currentSlot := time.Now().Unix() / 12
	recentSlot := currentSlot - (60 * 5) // 5 minutes ago

	tests := []struct {
		name    string
		slot    int64
		wantErr bool
	}{
		{
			name:    "Recent valid slot",
			slot:    recentSlot,
			wantErr: false,
		},
		{
			name:    "Future slot",
			slot:    currentSlot + 1000,
			wantErr: true,
		},
		{
			name:    "Very old slot",
			slot:    1,
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

				duties, err = service.GetSyncDutiesBySlot(ctx, tt.slot)
				if err == nil || tt.wantErr {
					break
				}

				// If it's not a rate limit error, don't retry
				if err != nil && err.Error() != "failed to get block: 429 Too Many Requests" {
					break
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSyncDutiesBySlot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Log the number of validators found
				t.Logf("Slot %d: Found %d validator duties", tt.slot, len(duties))

				// Verify that we got some validator public keys
				for i, pubKey := range duties {
					if len(pubKey) == 0 {
						t.Errorf("Empty public key found at index %d", i)
					}
				}
			}
		})
	}
}
