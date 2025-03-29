package handler

// BlockRewardResponse represents the response for a block reward query
type BlockRewardResponse struct {
	Status string `json:"status" example:"mev"` // mev or vanilla
	Reward int64  `json:"reward" example:"123456"` // reward in GWEI
	BlockInfo struct {
		ProposerPayment int64 `json:"proposer_payment" example:"123456"`
		IsMevBoost      bool  `json:"is_mev_boost" example:"true"`
	} `json:"block_info"`
}

// SyncDutiesResponse represents the response for a sync duties query
type SyncDutiesResponse struct {
	Validators []string `json:"validators" example:"['0x1234...','0x5678...']"`
	SyncInfo   struct {
		SyncPeriod    int64 `json:"sync_period" example:"123"`
		CommitteeSize int   `json:"committee_size" example:"32"`
	} `json:"sync_info"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
} 