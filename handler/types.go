package handler

// BlockRewardResponse represents the response structure for block rewards
type BlockRewardResponse struct {
	Status    string `json:"status" example:"mev" description:"mev or vanilla"` // Block type (MEV or vanilla)
	Reward    int64  `json:"reward" example:"123456" description:"reward in GWEI"` // Total block reward in GWEI
	BlockInfo struct {
		ProposerPayment int64 `json:"proposer_payment" example:"123456"` // Payment to block proposer in GWEI
		IsMEVBoost      bool  `json:"is_mev_boost" example:"true"`      // Whether MEV-Boost was used
	} `json:"block_info"`
}

// SyncDutiesResponse represents the response structure for sync committee duties
type SyncDutiesResponse struct {
	Validators []string `json:"validators" example:"['0x1234...','0x5678...']"` // List of validator public keys in the sync committee
	SyncInfo   struct {
		SyncPeriod    int64 `json:"sync_period" example:"123"`    // Current sync committee period number
		CommitteeSize int   `json:"committee_size" example:"512"` // Size of the sync committee
	} `json:"sync_info"`
}

// ErrorResponse represents the standard error response structure
type ErrorResponse struct {
	Error string `json:"error" example:"Internal server error"` // Error message
} 