package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Get Block Reward
// @Description Get the block reward for a given slot
// @Tags block
// @Param slot path int true "Slot Number"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /blockreward/{slot} [get]
func (h *Handler) GetBlockReward(c *gin.Context) {
	slotParam := c.Param("slot")
	slot, err := strconv.Atoi(slotParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot number"})
		return
	}

	// Placeholder logic for retrieving block reward
	// Replace with actual logic to fetch block reward data
	if slot < 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Slot does not exist"})
		return
	}

	if slot > 1000000 { // Example future slot check
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slot is in the future"})
		return
	}

	// Example response
	c.JSON(http.StatusOK, gin.H{
		"status": "vanilla",
		"reward": 1000, // Example reward in GWEI
	})
}
