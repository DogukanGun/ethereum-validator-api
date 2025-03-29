package handler

import "ethereum-validator-api/service"

// Handler manages HTTP request handling and coordinates with the Ethereum service
type Handler struct {
	ethService *service.EthereumService
}

// NewHandler creates a new Handler instance with the provided Ethereum service
func NewHandler(ethService *service.EthereumService) *Handler {
	return &Handler{
		ethService: ethService,
	}
}
