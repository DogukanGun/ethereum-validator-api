package handler

import (
	"errors"
	"ethereum-validator-api/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Get Sync Duties
// @Description Get the sync committee duties for validators at a given slot in the PoS chain
// @Tags sync
// @Param slot path int true "Slot Number"
// @Success 200 {object} SyncDutiesResponse "Returns list of validator public keys with sync committee duties"
// @Failure 400 {object} ErrorResponse "Invalid slot number or slot too far in future"
// @Failure 404 {object} ErrorResponse "Slot does not exist"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /syncduties/{slot} [get]
func (h *Handler) GetSyncDuties(c *gin.Context) {
	slotParam := c.Param("slot")
	slot, err := strconv.ParseInt(slotParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid slot number"})
		return
	}

	validators, err := h.ethService.GetSyncDutiesBySlot(c.Request.Context(), slot)
	if err != nil {
		var statusCode int
		var errMsg string

		switch {
		case errors.Is(err, service.ErrFutureSlot):
			statusCode = http.StatusBadRequest
			errMsg = "Slot is too far in the future"
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

	// Calculate sync period
	syncPeriod := slot / 8192 // Sync committee period changes every 256 epochs (8192 slots)

	// Create response object
	response := SyncDutiesResponse{
		Validators: validators,
	}
	response.SyncInfo.SyncPeriod = syncPeriod
	response.SyncInfo.CommitteeSize = len(validators)

	c.JSON(http.StatusOK, response)
}
