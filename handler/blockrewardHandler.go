package handler

import (
	"errors"
	"ethereum-validator-api/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Get Block Reward
// @Description Get the block reward and MEV information for a given slot
// @Tags block
// @Param slot path int true "Slot Number"
// @Success 200 {object} BlockRewardResponse "Returns block reward info including MEV status and reward in GWEI"
// @Failure 400 {object} ErrorResponse "Invalid slot number or future slot"
// @Failure 404 {object} ErrorResponse "Slot does not exist"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /blockreward/{slot} [get]
func (h *Handler) GetBlockReward(c *gin.Context) {
	slotParam := c.Param("slot")
	slot, err := strconv.ParseInt(slotParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid slot number"})
		return
	}

	reward, err := h.ethService.GetBlockRewardBySlot(c.Request.Context(), slot)
	if err != nil {
		var statusCode int
		var errMsg string

		switch {
		case errors.Is(err, service.ErrFutureSlot):
			statusCode = http.StatusBadRequest
			errMsg = "Slot is in the future"
		case errors.Is(err, service.ErrSlotNotFound):
			statusCode = http.StatusNotFound
			errMsg = "Slot does not exist"
		default:
			statusCode = http.StatusInternalServerError
			errMsg = "Internal server error"
		}

		c.JSON(statusCode, ErrorResponse{Error: errMsg})
		return
	}

	// Create response object
	response := BlockRewardResponse{
		Status: reward.Status,
		Reward: reward.Reward.Int64(),
	}
	response.BlockInfo.ProposerPayment = reward.Reward.Int64()
	response.BlockInfo.IsMevBoost = reward.Status == "mev"

	c.JSON(http.StatusOK, response)
}
