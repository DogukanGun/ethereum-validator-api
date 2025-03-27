package handler

import "ethereum-validator-api/service"

type Handler struct {
	ethService *service.EthereumService
}

func NewHandler(ethService *service.EthereumService) *Handler {
	return &Handler{
		ethService: ethService,
	}
}
