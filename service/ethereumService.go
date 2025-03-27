package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	_ "io"
	"math/big"
	"net/http"
	"net/url"
	_ "os"
	"strings"
	"time"
)

type EthereumService struct {
	rpcURL string
	client *http.Client
}

type BlockReward struct {
	Status string   `json:"status"` // "mev" or "vanilla"
	Reward *big.Int `json:"reward"` // in GWEI
}

// BeaconBlockResponse represents the response from the Beacon API for block details
type BeaconBlockResponse struct {
	Data struct {
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
	} `json:"data"`
}

// ExecutionBlockResponse represents the response from the Execution API
type ExecutionBlockResponse struct {
	Result struct {
		Transactions []struct {
			Hash             string `json:"hash"`
			GasPrice         string `json:"gasPrice"`
			GasUsed          string `json:"gasUsed"`
			MaxPriorityFee   string `json:"maxPriorityFeePerGas"`
			MaxFeePerGas     string `json:"maxFeePerGas"`
			TransactionIndex string `json:"transactionIndex"`
		} `json:"transactions"`
		BaseFeePerGas string `json:"baseFeePerGas"`
	} `json:"result"`
}

// SyncCommitteeResponse represents the response from the Beacon API for sync committee duties
type SyncCommitteeResponse struct {
	Data struct {
		ValidatorSyncAssignments []struct {
			ValidatorPubKey string `json:"validator_pubkey"`
		} `json:"validator_sync_assignments"`
	} `json:"data"`
}

// RPCRequest represents a JSON-RPC request
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// Known MEV-Boost builder prefixes in extraData
var mevBuilderPrefixes = []string{
	"flashbots",
	"builder0x69",
	"rsync-builder",
	"manifold",
	"eth-builder",
}

func NewEthereumService(rpcURL string) (*EthereumService, error) {
	// Validate URL
	if rpcURL == "" {
		return nil, fmt.Errorf("RPC URL cannot be empty")
	}

	parsedURL, err := url.Parse(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("invalid RPC URL: %v", err)
	}

	// Additional URL validation
	if !parsedURL.IsAbs() {
		return nil, fmt.Errorf("RPC URL must be absolute")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("RPC URL must use http or https scheme")
	}

	return &EthereumService{
		rpcURL: rpcURL,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

// GetBlockRewardBySlot retrieves block reward information for a given slot
func (s *EthereumService) GetBlockRewardBySlot(ctx context.Context, slot int64) (*BlockReward, error) {
	// Validate slot is not in the future
	currentSlot := time.Now().Unix() / 12 // 12 second slots
	if slot > currentSlot {
		return nil, fmt.Errorf("requested slot %d is in the future (current slot: %d)", slot, currentSlot)
	}

	// First get the beacon block to check if it's MEV
	beaconBlock, err := s.getBeaconBlock(ctx, slot)
	if err != nil {
		return nil, fmt.Errorf("failed to get beacon block: %v", err)
	}

	// Check if block is MEV produced
	isMev := s.isMEVBlock(beaconBlock)

	// Get execution block details for reward calculation
	blockHash := beaconBlock.Data.Message.Body.ExecutionPayload.BlockHash
	reward, err := s.getExecutionBlockReward(ctx, blockHash, beaconBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution block reward: %v", err)
	}

	// Convert Wei to Gwei
	gweiReward := new(big.Int).Div(reward, big.NewInt(1e9))

	return &BlockReward{
		Status: map[bool]string{true: "mev", false: "vanilla"}[isMev],
		Reward: gweiReward,
	}, nil
}

// isMEVBlock checks if a block was produced by MEV-Boost
func (s *EthereumService) isMEVBlock(block *BeaconBlockResponse) bool {
	extraData := block.Data.Message.Body.ExecutionPayload.ExtraData

	// Check for empty extraData
	if len(extraData) == 0 {
		return false
	}

	// Check for known MEV builder signatures in extraData
	for _, prefix := range mevBuilderPrefixes {
		if strings.Contains(strings.ToLower(extraData), prefix) {
			return true
		}
	}

	// Additional MEV detection logic based on transaction patterns
	// MEV blocks often have specific transaction ordering and patterns
	txs := block.Data.Message.Body.ExecutionPayload.Transactions
	if len(txs) > 0 {
		// Check for sandwich patterns (common in MEV blocks)
		// This is a simplified check - real MEV detection would be more sophisticated
		if len(txs) >= 3 && txs[0] == txs[len(txs)-1] {
			return true
		}
	}

	return false
}

// GetSyncDutiesBySlot retrieves sync committee duties for a given slot
func (s *EthereumService) GetSyncDutiesBySlot(ctx context.Context, slot int64) ([]string, error) {
	// Validate slot
	currentSlot := time.Now().Unix() / 12 // 12 second slots
	if slot > currentSlot {
		return nil, fmt.Errorf("requested slot is in the future")
	}

	// Calculate sync committee period
	// Sync committee period changes every 256 epochs (8192 slots)
	syncPeriod := slot / 8192

	// Prepare JSON-RPC request
	rpcReq := RPCRequest{
		JSONRPC: "2.0",
		Method:  "beacon_get_sync_committee",
		Params:  []interface{}{fmt.Sprintf("0x%x", syncPeriod)},
		ID:      1,
	}

	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.rpcURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var syncResponse SyncCommitteeResponse
	if err := json.NewDecoder(resp.Body).Decode(&syncResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var validatorPubKeys []string
	for _, assignment := range syncResponse.Data.ValidatorSyncAssignments {
		validatorPubKeys = append(validatorPubKeys, assignment.ValidatorPubKey)
	}

	return validatorPubKeys, nil
}

func (s *EthereumService) getBeaconBlock(ctx context.Context, slot int64) (*BeaconBlockResponse, error) {
	rpcReq := RPCRequest{
		JSONRPC: "2.0",
		Method:  "beacon_get_block",
		Params:  []interface{}{fmt.Sprintf("0x%x", slot)},
		ID:      1,
	}

	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.rpcURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var blockResponse BeaconBlockResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockResponse); err != nil {
		return nil, err
	}

	return &blockResponse, nil
}

func (s *EthereumService) getExecutionBlockReward(ctx context.Context, blockHash string, beaconBlock *BeaconBlockResponse) (*big.Int, error) {
	rpcReq := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByHash",
		Params:  []interface{}{blockHash, true},
		ID:      1,
	}

	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.rpcURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var blockResponse ExecutionBlockResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockResponse); err != nil {
		return nil, err
	}

	totalReward := new(big.Int)

	// Safely parse base fee
	baseFeePerGas := new(big.Int)
	if blockResponse.Result.BaseFeePerGas != "" {
		baseFeeHex := strings.TrimPrefix(blockResponse.Result.BaseFeePerGas, "0x")
		if _, ok := baseFeePerGas.SetString(baseFeeHex, 16); !ok {
			return nil, fmt.Errorf("failed to parse base fee: %s", blockResponse.Result.BaseFeePerGas)
		}
	}

	// Calculate rewards for each transaction
	for _, tx := range blockResponse.Result.Transactions {
		// Calculate priority fee
		var priorityFee *big.Int
		if tx.MaxPriorityFee != "" {
			priorityFee = new(big.Int)
			priorityHex := strings.TrimPrefix(tx.MaxPriorityFee, "0x")
			if _, ok := priorityFee.SetString(priorityHex, 16); !ok {
				continue // Skip transaction if priority fee parsing fails
			}
		} else if tx.GasPrice != "" {
			// For legacy transactions, priority fee is gasPrice - baseFee
			gasPrice := new(big.Int)
			gasPriceHex := strings.TrimPrefix(tx.GasPrice, "0x")
			if _, ok := gasPrice.SetString(gasPriceHex, 16); !ok {
				continue // Skip transaction if gas price parsing fails
			}
			priorityFee = new(big.Int).Sub(gasPrice, baseFeePerGas)
			if priorityFee.Sign() < 0 {
				priorityFee = big.NewInt(0)
			}
		} else {
			continue // Skip if neither priority fee nor gas price is available
		}

		// Parse gas used
		gasUsed := new(big.Int)
		gasUsedHex := strings.TrimPrefix(tx.GasUsed, "0x")
		if _, ok := gasUsed.SetString(gasUsedHex, 16); !ok {
			continue // Skip transaction if gas used parsing fails
		}

		// Calculate transaction reward (priority fee * gas used)
		txReward := new(big.Int).Mul(priorityFee, gasUsed)
		totalReward.Add(totalReward, txReward)
	}

	return totalReward, nil
}
