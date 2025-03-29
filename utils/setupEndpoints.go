package utils

import (
	"ethereum-validator-api/handler"
	"ethereum-validator-api/service"
	"github.com/gin-gonic/gin"
	"os"
)

// SetupEndpoints configures the API endpoints for the Ethereum validator service
func SetupEndpoints(router *gin.Engine) error {
	rpcURL := os.Getenv("ETH_RPC")
	ethService, err := service.NewEthereumService(rpcURL)
	if err != nil {
		return err
	}

	h := handler.NewHandler(ethService)

	// Register API endpoints
	router.GET("/blockreward/:slot", h.GetBlockReward)
	router.GET("/syncduties/:slot", h.GetSyncDuties)

	return nil
}
