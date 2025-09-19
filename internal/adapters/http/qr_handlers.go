package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"haoma/internal/config"
)

// ScanNodeQR godoc
// @Summary Scan QR code to access a carnival node
// @Description Scan a QR code at a physical location to unlock and start a carnival node
// @Tags Nodes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body StartNodeRequest true "QR code scan information"
// @Success 200 {object} StartNodeResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /nodes/scan [post]
func (h *CarnivalHandler) ScanNodeQR(c *gin.Context) {
	var req StartNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	playerID, exists := c.Get("player_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Player not authenticated"})
		return
	}

	resultSessionID, node, err := h.service.ScanNodeQR(playerID.(uuid.UUID), req.NodeCode, req.SessionID)
	if err != nil {
		if err.Error() == "node not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid QR code - node not found"})
			return
		}
		if err.Error() == "node already completed" {
			c.JSON(http.StatusConflict, gin.H{"error": "You have already completed this node"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	nodeResp := NodeResponse{
		Number:    node.Number,
		Questions: make([]QuestionResponse, len(node.Questions)),
	}

	for i, q := range node.Questions {
		nodeResp.Questions[i] = QuestionResponse{
			ID:      q.ID,
			Text:    q.Text,
			OptionA: q.OptionA,
			OptionB: q.OptionB,
			OptionC: q.OptionC,
			OptionD: q.OptionD,
		}
	}

	c.JSON(http.StatusOK, StartNodeResponse{
		SessionID: *resultSessionID,
		Node:      nodeResp,
		Message:   fmt.Sprintf(config.WELCOME_NODE_MESSAGE, node.Number, config.QUESTIONS_PER_NODE),
	})
}
