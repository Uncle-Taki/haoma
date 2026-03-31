package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetInventory returns inventory items, optionally filtered by category.
// @Summary Get inventory items
// @Description Retrieve available medicines, food, and supplements
// @Tags Inventory
// @Produce json
// @Param category query string false "Filter by category: medicine, food, supplement"
// @Success 200 {array} inventory.InventoryItem
// @Failure 400 {object} map[string]interface{}
// @Router /vet/inventory [get]
func (h *VetHandler) GetInventory(c *gin.Context) {
	category := c.Query("category")

	items, err := h.service.GetInventory(category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetInventoryItem returns a single inventory item by ID.
// @Summary Get inventory item details
// @Description Retrieve details for a specific inventory item
// @Tags Inventory
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} inventory.InventoryItem
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vet/inventory/:id [get]
func (h *VetHandler) GetInventoryItem(c *gin.Context) {
	idParam := c.Param("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.service.GetInventoryItem(itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}
