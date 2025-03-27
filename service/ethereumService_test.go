package service

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewEthereumService(t *testing.T) {
	tests := []struct {
		name        string
		rpcURL      string
		wantErr     bool
		errContains string
	}{
		{
			name:    "Valid URL",
			rpcURL:  "https://example.com",
			wantErr: false,
		},
		{
			name:        "Empty URL",
			rpcURL:      "",
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "Invalid URL",
			rpcURL:      "not-a-url",
			wantErr:     true,
			errContains: "invalid RPC URL",
		},
		{
			name:        "Non-absolute URL",
			rpcURL:      "path/to/somewhere",
			wantErr:     true,
			errContains: "must be absolute",
		},
		{
			name:        "Invalid scheme",
			rpcURL:      "ftp://example.com",
			wantErr:     true,
			errContains: "must use http or https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEthereumService(tt.rpcURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEthereumService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Error("NewEthereumService() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("NewEthereumService() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewEthereumService() returned nil service")
			}
		})
	}
}

func TestEthereumService_GetBlockRewardBySlot(t *testing.T) {
	// Calculate current slot for test cases
	currentSlot := time.Now().Unix() / 12
	futureSlot := currentSlot + 1000
	recentSlot := currentSlot - 100
	oldSlot := currentSlot - 10000

	tests := []struct {
		name           string
		slot          int64
		beaconResp    BeaconBlockResponse
		executionResp ExecutionBlockResponse
		wantStatus    string
		wantReward    *big.Int
		wantErr       bool
		errorContains string
	}{
		{
			name:        "Future slot",
			slot:        futureSlot,
			wantErr:     true,
			errorContains: "is in the future",
		},
		{
			name: "Recent valid slot",
			slot: recentSlot,
			beaconResp: BeaconBlockResponse{
				Data: struct {
					Message struct {
						Body struct {
							ExecutionPayload struct {
								FeeRecipient  string   `json:"fee_recipient"`
								BlockHash     string   `json:"block_hash"`
								ExtraData     string   `json:"extra_data"`
								Transactions  []string `json:"transactions"`
								BaseFeePerGas string   `json:"base_fee_per_gas"`
							} `json:"execution_payload"`
						} `json:"body"`
						ProposerIndex string `json:"proposer_index"`
					} `json:"message"`
				}{
					Message: struct {
						Body struct {
							ExecutionPayload struct {
								FeeRecipient  string   `json:"fee_recipient"`
								BlockHash     string   `json:"block_hash"`
								ExtraData     string   `json:"extra_data"`
								Transactions  []string `json:"transactions"`
								BaseFeePerGas string   `json:"base_fee_per_gas"`
							} `json:"execution_payload"`
						} `json:"body"`
						ProposerIndex string `json:"proposer_index"`
					}{
						Body: struct {
							ExecutionPayload struct {
								FeeRecipient  string   `json:"fee_recipient"`
								BlockHash     string   `json:"block_hash"`
								ExtraData     string   `json:"extra_data"`
								Transactions  []string `json:"transactions"`
								BaseFeePerGas string   `json:"base_fee_per_gas"`
							} `json:"execution_payload"`
						}{
							ExecutionPayload: struct {
								FeeRecipient  string   `json:"fee_recipient"`
								BlockHash     string   `json:"block_hash"`
								ExtraData     string   `json:"extra_data"`
								Transactions  []string `json:"transactions"`
								BaseFeePerGas string   `json:"base_fee_per_gas"`
							}{
								ExtraData:     "",
								BlockHash:     "0x123",
								BaseFeePerGas: "0x5",
								Transactions:  []string{"0x1"},
							},
						},
					},
				},
			},
			executionResp: ExecutionBlockResponse{
				Result: struct {
					Transactions []struct {
						Hash             string `json:"hash"`
						GasPrice         string `json:"gasPrice"`
						GasUsed          string `json:"gasUsed"`
						MaxPriorityFee   string `json:"maxPriorityFeePerGas"`
						MaxFeePerGas     string `json:"maxFeePerGas"`
						TransactionIndex string `json:"transactionIndex"`
					} `json:"transactions"`
					BaseFeePerGas string `json:"baseFeePerGas"`
				}{
					Transactions: []struct {
						Hash             string `json:"hash"`
						GasPrice         string `json:"gasPrice"`
						GasUsed          string `json:"gasUsed"`
						MaxPriorityFee   string `json:"maxPriorityFeePerGas"`
						MaxFeePerGas     string `json:"maxFeePerGas"`
						TransactionIndex string `json:"transactionIndex"`
					}{
						{
							GasPrice: "0x8",
							GasUsed:  "0x5208",
						},
					},
					BaseFeePerGas: "0x5",
				},
			},
			wantStatus: "vanilla",
			wantReward: new(big.Int).Mul(big.NewInt(3), big.NewInt(21000)), // (gasPrice - baseFee) * gasUsed
			wantErr:    false,
		},
		{
			name: "Very old slot",
			slot: oldSlot,
			beaconResp: BeaconBlockResponse{
				Data: struct {
					Message struct {
						Body struct {
							ExecutionPayload struct {
								FeeRecipient  string   `json:"fee_recipient"`
								BlockHash     string   `json:"block_hash"`
								ExtraData     string   `json:"extra_data"`
								Transactions  []string `json:"transactions"`
								BaseFeePerGas string   `json:"base_fee_per_gas"`
							} `json:"execution_payload"`
						} `json:"body"`
						ProposerIndex string `json:"proposer_index"`
					} `json:"message"`
				}{
					Message: struct {
						Body struct {
							ExecutionPayload struct {
								FeeRecipient  string   `json:"fee_recipient"`
								BlockHash     string   `json:"block_hash"`
								ExtraData     string   `json:"extra_data"`
								Transactions  []string `json:"transactions"`
								BaseFeePerGas string   `json:"base_fee_per_gas"`
							} `json:"execution_payload"`
						} `json:"body"`
						ProposerIndex string `json:"proposer_index"`
					}{
						Body: struct {
							ExecutionPayload struct {
								FeeRecipient  string   `json:"fee_recipient"`
								BlockHash     string   `json:"block_hash"`
								ExtraData     string   `json:"extra_data"`
								Transactions  []string `json:"transactions"`
								BaseFeePerGas string   `json:"base_fee_per_gas"`
							} `json:"execution_payload"`
						}{
							ExecutionPayload: struct {
								FeeRecipient  string   `json:"fee_recipient"`
								BlockHash     string   `json:"block_hash"`
								ExtraData     string   `json:"extra_data"`
								Transactions  []string `json:"transactions"`
								BaseFeePerGas string   `json:"base_fee_per_gas"`
							}{
								ExtraData:     "",
								BlockHash:     "0x456",
								BaseFeePerGas: "0x5",
								Transactions:  []string{},
							},
						},
					},
				},
			},
			executionResp: ExecutionBlockResponse{
				Result: struct {
					Transactions []struct {
						Hash             string `json:"hash"`
						GasPrice         string `json:"gasPrice"`
						GasUsed          string `json:"gasUsed"`
						MaxPriorityFee   string `json:"maxPriorityFeePerGas"`
						MaxFeePerGas     string `json:"maxFeePerGas"`
						TransactionIndex string `json:"transactionIndex"`
					} `json:"transactions"`
					BaseFeePerGas string `json:"baseFeePerGas"`
				}{
					Transactions: []struct {
						Hash             string `json:"hash"`
						GasPrice         string `json:"gasPrice"`
						GasUsed          string `json:"gasUsed"`
						MaxPriorityFee   string `json:"maxPriorityFeePerGas"`
						MaxFeePerGas     string `json:"maxFeePerGas"`
						TransactionIndex string `json:"transactionIndex"`
					}{},
					BaseFeePerGas: "0x5",
				},
			},
			wantStatus: "vanilla",
			wantReward: big.NewInt(0), // Empty block
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server with mock responses
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req RPCRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatalf("Failed to decode request: %v", err)
				}

				switch req.Method {
				case "beacon_get_block":
					json.NewEncoder(w).Encode(tt.beaconResp)
				case "eth_getBlockByHash":
					json.NewEncoder(w).Encode(tt.executionResp)
				default:
					t.Fatalf("Unexpected method: %s", req.Method)
				}
			}))
			defer server.Close()

			// Create service with test server URL
			s := &EthereumService{
				rpcURL: server.URL,
				client: server.Client(),
			}

			got, err := s.GetBlockRewardBySlot(context.Background(), tt.slot)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockRewardBySlot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errorContains != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("GetBlockRewardBySlot() error = %v, want error containing %v", err, tt.errorContains)
				}
				return
			}
			if err == nil {
				if got.Status != tt.wantStatus {
					t.Errorf("GetBlockRewardBySlot() status = %v, want %v", got.Status, tt.wantStatus)
				}
				if got.Reward.Cmp(tt.wantReward) != 0 {
					t.Errorf("GetBlockRewardBySlot() reward = %v, want %v", got.Reward, tt.wantReward)
				}
			}
		})
	}
}

func TestEthereumService_GetSyncDutiesBySlot(t *testing.T) {
	tests := []struct {
		name         string
		slot         int64
		syncResponse SyncCommitteeResponse
		wantKeys     []string
		wantErr      bool
	}{
		{
			name: "Valid sync committee response",
			slot: 1000,
			syncResponse: SyncCommitteeResponse{
				Data: struct {
					ValidatorSyncAssignments []struct {
						ValidatorPubKey string `json:"validator_pubkey"`
					} `json:"validator_sync_assignments"`
				}{
					ValidatorSyncAssignments: []struct {
						ValidatorPubKey string `json:"validator_pubkey"`
					}{
						{ValidatorPubKey: "0x123"},
						{ValidatorPubKey: "0x456"},
					},
				},
			},
			wantKeys: []string{"0x123", "0x456"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(tt.syncResponse)
			}))
			defer server.Close()

			s := &EthereumService{
				rpcURL: server.URL,
				client: server.Client(),
			}

			got, err := s.GetSyncDutiesBySlot(context.Background(), tt.slot)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSyncDutiesBySlot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if len(got) != len(tt.wantKeys) {
					t.Errorf("GetSyncDutiesBySlot() got %v keys, want %v", len(got), len(tt.wantKeys))
				}
				for i, key := range got {
					if key != tt.wantKeys[i] {
						t.Errorf("GetSyncDutiesBySlot() key[%d] = %v, want %v", i, key, tt.wantKeys[i])
					}
				}
			}
		})
	}
} 