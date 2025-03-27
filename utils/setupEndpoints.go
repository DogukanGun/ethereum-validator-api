package utils

import (
	"ethereum-validator-api/handler"
	"ethereum-validator-api/service"
	"github.com/gin-gonic/gin"
	"os"
)

// SetupEndpoints configures and initializes all the API endpoints for the Ethereum validator service.
// It creates a new Ethereum service instance and sets up the HTTP routes.
//
// Parameters:
//   - router: A pointer to the Gin router engine to register the routes
//
// Returns:
//   - error: Returns an error if there's any issue initializing the Ethereum service,
//     nil otherwise
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
